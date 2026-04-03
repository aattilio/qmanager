package hypervisor

import (
	"fmt"

	"github.com/libvirt/libvirt-go"
)

type NetworkInfrastructureManager struct {
	LibvirtConnection *libvirt.Connect
}

func NewNetworkInfrastructureManager(
	connection *libvirt.Connect,
) *NetworkInfrastructureManager {
	return &NetworkInfrastructureManager{
		LibvirtConnection: connection,
	}
}

func (manager *NetworkInfrastructureManager) CreateBridgeNetwork(
	networkName string,
	bridgeDeviceName string,
) error {
	xmlConfiguration := fmt.Sprintf(
		`<network>
		  <name>%s</name>
		  <forward mode='bridge'/>
		  <bridge name='%s'/>
		</network>`,
		networkName,
		bridgeDeviceName,
	)
	
	networkDefinition, err := manager.LibvirtConnection.NetworkDefineXML(xmlConfiguration)
	if err != nil {
		return err
	}
	defer networkDefinition.Free()
	
	return networkDefinition.Create()
}

func (manager *NetworkInfrastructureManager) CreateNatNetwork(
	networkName string,
	bridgeDeviceName string,
	ipv4AddressRange string,
) error {
	xmlConfiguration := fmt.Sprintf(
		`<network>
		  <name>%s</name>
		  <forward mode='nat'/>
		  <bridge name='%s' stp='on' delay='0'/>
		  <ip address='%s' netmask='255.255.255.0'>
		    <dhcp>
		      <range start='%s.2' end='%s.254'/>
		    </dhcp>
		  </ip>
		</network>`,
		networkName,
		bridgeDeviceName,
		ipv4AddressRange,
		ipv4AddressRange,
		ipv4AddressRange,
	)
	
	networkDefinition, err := manager.LibvirtConnection.NetworkDefineXML(xmlConfiguration)
	if err != nil {
		return err
	}
	defer networkDefinition.Free()
	
	return networkDefinition.Create()
}

func (manager *NetworkInfrastructureManager) GetActiveNetworkNames() ([]string, error) {
	networks, err := manager.LibvirtConnection.ListAllNetworks(
		libvirt.CONNECT_LIST_NETWORKS_ACTIVE,
	)
	if err != nil {
		return nil, err
	}
	defer func() {
		for _, network := range networks {
			network.Free()
		}
	}()
	
	names := make([]string, 0, len(networks))
	for _, network := range networks {
		name, err := network.GetName()
		if err == nil {
			names = append(
				names,
				name,
			)
		}
	}
	
	return names, nil
}
