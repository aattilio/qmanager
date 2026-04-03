package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"

	"qmanager/src/backend/hypervisor"
	"qmanager/src/backend/filesystem"
	"qmanager/src/core"
)

type AuditReport struct {
	HypervisorAvailable bool     `json:"hypervisor_available"`
	ActiveVirtualMachines []string `json:"active_virtual_machines"`
	CatalogSystemsCount int      `json:"catalog_systems_count"`
	StoragePath         string   `json:"storage_path"`
}

func main() {
	action := flag.String(
		"action", 
		"audit", 
		"Action to perform (audit, list, provision, start, stop)",
	)
	
	vmName := flag.String(
		"name", 
		"", 
		"Virtual machine name",
	)
	
	osId := flag.String(
		"os", 
		"", 
		"Operating system ID for provisioning",
	)
	
	flag.Parse()

	fmt.Println("QManager - Advanced Hypervisor Control Interface")

	switch *action {
	case "audit":
		performInfrastructureAudit()
	case "list":
		listVirtualMachines()
	case "provision":
		provisionNewMachine(*osId, *vmName)
	case "start":
		controlMachine(*vmName, true)
	case "stop":
		controlMachine(*vmName, false)
	default:
		log.Fatalf("unsupported_action: %s", *action)
	}
}

func performInfrastructureAudit() {
	report := AuditReport{
		StoragePath: "data/vms",
	}

	connector, err := hypervisor.NewLibvirtHypervisorConnector("qemu:///system")
	if err == nil {
		report.HypervisorAvailable = true
		names, _ := connector.ListAllVirtualMachineNames()
		report.ActiveVirtualMachines = names
		connector.CloseConnection()
	}

	catalog, err := core.LoadConfigurationCatalogFromDirectory("config/catalog")
	if err == nil {
		report.CatalogSystemsCount = len(catalog.OperatingSystems)
	}

	output, _ := json.MarshalIndent(report, "", "  ")
	fmt.Println("Infrastructure Audit Report:")
	fmt.Println(string(output))
}

func listVirtualMachines() {
	connector, err := hypervisor.NewLibvirtHypervisorConnector("qemu:///system")
	if err != nil {
		log.Fatalf("hypervisor_unreachable: %v", err)
	}
	defer connector.CloseConnection()

	names, err := connector.ListAllVirtualMachineNames()
	if err != nil {
		log.Fatalf("failed_to_list_vms: %v", err)
	}

	fmt.Printf("Detected %d virtual machines:\n", len(names))
	for _, name := range names {
		fmt.Printf(" - %s\n", name)
	}
}

func provisionNewMachine(osId, vmName string) {
	if osId == "" || vmName == "" {
		log.Fatal("provisioning_requires_os_id_and_vm_name")
	}

	connector, err := hypervisor.NewLibvirtHypervisorConnector("qemu:///system")
	if err != nil {
		log.Fatal(err)
	}
	defer connector.CloseConnection()

	diskManager, _ := filesystem.NewVirtualDiskManager("data/vms")
	catalog, _ := core.LoadConfigurationCatalogFromDirectory("config/catalog")

	provisioner := core.NewAutomatedVirtualMachineProvisioner(
		connector, 
		diskManager, 
		"data",
	)

	fmt.Printf("Starting automated provisioning for %s...\n", vmName)
	err = provisioner.ExecuteExpressInstallation(osId, vmName, catalog)
	if err != nil {
		log.Fatalf("provisioning_failed: %v", err)
	}
	fmt.Println("Provisioning completed successfully.")
}

func controlMachine(name string, start bool) {
	if name == "" {
		log.Fatal("machine_name_required")
	}

	connector, err := hypervisor.NewLibvirtHypervisorConnector("qemu:///system")
	if err != nil {
		log.Fatal(err)
	}
	defer connector.CloseConnection()

	if start {
		err = connector.StartVirtualMachine(name)
	} else {
		err = connector.StopVirtualMachine(name)
	}

	if err != nil {
		log.Fatalf("operation_failed: %v", err)
	}
	fmt.Printf("Machine '%s' control signal sent successfully.\n", name)
}
