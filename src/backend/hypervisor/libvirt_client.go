package hypervisor

import (
	"fmt"

	"github.com/libvirt/libvirt-go"
)

type LibvirtHypervisorConnector struct {
	ActiveConnection *libvirt.Connect
}

func NewLibvirtHypervisorConnector(
	connectionUri string,
) (
	*LibvirtHypervisorConnector, 
	error,
) {
	connection, err := libvirt.NewConnect(
		connectionUri,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"hypervisor_connection_failed: %w", 
			err,
		)
	}
	
	return &LibvirtHypervisorConnector{
		ActiveConnection: connection,
	}, nil
}

func (connector *LibvirtHypervisorConnector) CloseConnection() error {
	if connector.ActiveConnection != nil {
		_, err := connector.ActiveConnection.Close()
		return err
	}
	
	return nil
}

func (connector *LibvirtHypervisorConnector) ListAllVirtualMachineNames() (
	[]string, 
	error,
) {
	domains, err := connector.ActiveConnection.ListAllDomains(
		libvirt.CONNECT_LIST_DOMAINS_ACTIVE | libvirt.CONNECT_LIST_DOMAINS_INACTIVE,
	)
	if err != nil {
		return nil, err
	}
	
	defer func() {
		for _, domain := range domains {
			domain.Free()
		}
	}()
	
	names := make(
		[]string, 
		0, 
		len(domains),
	)
	
	for _, domain := range domains {
		name, err := domain.GetName()
		if err == nil {
			names = append(
				names, 
				name,
			)
		}
	}
	
	return names, nil
}

func (connector *LibvirtHypervisorConnector) StartVirtualMachine(
	machineName string,
) error {
	domain, err := connector.ActiveConnection.LookupDomainByName(
		machineName,
	)
	if err != nil {
		return err
	}
	defer domain.Free()
	
	return domain.Create()
}

func (connector *LibvirtHypervisorConnector) StopVirtualMachine(
	machineName string,
) error {
	domain, err := connector.ActiveConnection.LookupDomainByName(
		machineName,
	)
	if err != nil {
		return err
	}
	defer domain.Free()
	
	return domain.Destroy()
}

func (connector *LibvirtHypervisorConnector) DefineVirtualMachine(
	xmlDefinition string,
) error {
	domain, err := connector.ActiveConnection.DomainDefineXML(
		xmlDefinition,
	)
	if err != nil {
		return err
	}
	defer domain.Free()
	
	return nil
}

func (connector *LibvirtHypervisorConnector) CreateSnapshot(
	machineName string,
	snapshotName string,
) error {
	domain, err := connector.ActiveConnection.LookupDomainByName(
		machineName,
	)
	if err != nil {
		return err
	}
	defer domain.Free()

	xml := fmt.Sprintf(
		"<domainsnapshot><name>%s</name></domainsnapshot>",
		snapshotName,
	)
	
	snapshot, err := domain.CreateSnapshotXML(
		xml,
		libvirt.DOMAIN_SNAPSHOT_CREATE_ATOMIC,
	)
	if err != nil {
		return err
	}
	defer snapshot.Free()
	
	return nil
}

func (connector *LibvirtHypervisorConnector) GetMachineStats(
	machineName string,
) (
	uint64, 
	uint64, 
	error,
) {
	domain, err := connector.ActiveConnection.LookupDomainByName(
		machineName,
	)
	if err != nil {
		return 0, 0, err
	}
	defer domain.Free()

	stats, err := domain.MemoryStats(
		uint32(libvirt.DOMAIN_MEMORY_STAT_NR),
		0,
	)
	if err != nil {
		return 0, 0, err
	}

	var rss uint64
	for _, stat := range stats {
		if stat.Tag == int32(libvirt.DOMAIN_MEMORY_STAT_RSS) {
			rss = stat.Val
		}
	}

	cpuStats, err := domain.GetCPUStats(
		-1,
		1,
		0,
	)
	if err != nil || len(cpuStats) == 0 {
		return rss, 0, nil
	}

	return rss, cpuStats[0].CpuTime, nil
}
