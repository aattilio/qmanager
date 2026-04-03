package tests

import (
	"strings"
	"testing"

	"qmanager/src/backend/hypervisor"
)

func TestGenerateVirtualMachineXmlDefinition(t *testing.T) {
	machineName := "test-vm"
	ramMegabytes := 2048
	cpuCores := 4
	diskPath := "/tmp/disk.qcow2"
	isoPath := "/tmp/os.iso"

	xml, err := hypervisor.GenerateVirtualMachineXmlDefinition(
		machineName,
		ramMegabytes,
		cpuCores,
		diskPath,
		isoPath,
	)
	if err != nil {
		t.Fatal(err)
	}

	expectedSnippets := []string{
		"domain type=\"kvm\"",
		"<name>test-vm</name>",
		"memory unit=\"KiB\">2097152</memory>",
		"<vcpu>4</vcpu>",
		"driver name=\"qemu\" type=\"qcow2\"",
		"source file=\"/tmp/disk.qcow2\"",
		"device=\"cdrom\"",
		"source file=\"/tmp/os.iso\"",
		"mode=\"host-passthrough\"",
		"type=\"spice\"",
	}

	for _, snippet := range expectedSnippets {
		if !strings.Contains(xml, snippet) {
			t.Errorf("missing_xml_element: %s\nFull XML:\n%s", snippet, xml)
		}
	}
}
