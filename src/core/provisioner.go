package core

import (
	"fmt"
	"path/filepath"

	"qmanager/src/backend/hypervisor"
	"qmanager/src/backend/provisioning"
	"qmanager/src/backend/filesystem"
)

type AutomatedVirtualMachineProvisioner struct {
	LibvirtHypervisorConnector *hypervisor.LibvirtHypervisorConnector
	VirtualDiskManager         *filesystem.VirtualDiskManager
	BaseDataDirectory          string
}

func NewAutomatedVirtualMachineProvisioner(
	hypervisorConnector *hypervisor.LibvirtHypervisorConnector,
	diskManager *filesystem.VirtualDiskManager,
	dataDirectory string,
) *AutomatedVirtualMachineProvisioner {
	return &AutomatedVirtualMachineProvisioner{
		LibvirtHypervisorConnector: hypervisorConnector,
		VirtualDiskManager:         diskManager,
		BaseDataDirectory:          dataDirectory,
	}
}

func (provisioner *AutomatedVirtualMachineProvisioner) ExecuteExpressInstallation(
	operatingSystemId string,
	virtualMachineName string,
	configurationCatalog *Catalog,
) error {
	var selectedOperatingSystem *OperatingSystemMetadata
	for _, osMetadata := range configurationCatalog.OperatingSystems {
		if osMetadata.ID == operatingSystemId {
			selectedOperatingSystem = &osMetadata
			break
		}
	}

	if selectedOperatingSystem == nil {
		return fmt.Errorf(
			"operating_system_not_found_in_catalog: %s",
			operatingSystemId,
		)
	}

	isoStoragePath := filepath.Join(
		provisioner.BaseDataDirectory,
		"iso_cache",
		operatingSystemId+".iso",
	)
	
	downloadTask := provisioning.DownloadTask{
		URL:  selectedOperatingSystem.Mirrors[0],
		Dest: isoStoragePath,
		Progress: func(
			currentBytes, 
			totalBytes int64,
		) {
		},
	}

	if err := provisioning.ExecuteDownload(
		downloadTask,
	); err != nil {
		return fmt.Errorf(
			"failed_to_download_operating_system_iso: %w",
			err,
		)
	}

	diskSizeGigabytes := 40
	if selectedOperatingSystem.MinDiskGB > 0 {
		diskSizeGigabytes = selectedOperatingSystem.MinDiskGB
	}

	virtualDiskPath, err := provisioner.VirtualDiskManager.CreateQcow2(
		virtualMachineName,
		diskSizeGigabytes,
	)
	if err != nil {
		return fmt.Errorf(
			"failed_to_create_virtual_disk_image: %w",
			err,
		)
	}

	ramMegabytes := 4096
	if selectedOperatingSystem.MinRAM > 0 {
		ramMegabytes = selectedOperatingSystem.MinRAM
	}

	cpuCores := 2
	if selectedOperatingSystem.MinVCPUs > 0 {
		cpuCores = selectedOperatingSystem.MinVCPUs
	}

	xmlDefinition, err := hypervisor.GenerateVirtualMachineXmlDefinition(
		virtualMachineName,
		ramMegabytes,
		cpuCores,
		virtualDiskPath,
		isoStoragePath,
	)
	if err != nil {
		return fmt.Errorf(
			"failed_to_generate_virtual_machine_xml_definition: %w",
			err,
		)
	}

	if err := provisioner.LibvirtHypervisorConnector.DefineVirtualMachine(
		xmlDefinition,
	); err != nil {
		return fmt.Errorf(
			"failed_to_define_virtual_machine_in_libvirt: %w",
			err,
		)
	}

	return nil
}
