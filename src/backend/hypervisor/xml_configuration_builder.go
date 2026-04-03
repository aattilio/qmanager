package hypervisor

import (
	"encoding/xml"
	"fmt"
	"runtime"
)

type VirtualMachineConfiguration struct {
	XMLName             xml.Name            `xml:"domain"`
	HypervisorType      string              `xml:"type,attr"`
	MachineName         string              `xml:"name"`
	MemoryAllocation    MemoryAllocation    `xml:"memory"`
	CurrentMemory       MemoryAllocation    `xml:"currentMemory"`
	MemoryBacking       *MemoryBacking      `xml:"memoryBacking,omitempty"`
	VirtualProcessorCount int               `xml:"vcpu"`
	OperatingSystem     OperatingSystem     `xml:"os"`
	SystemFeatures      SystemFeatures      `xml:"features"`
	ProcessorModel      ProcessorModel      `xml:"cpu"`
	SystemClock         SystemClock         `xml:"clock"`
	HardwareDevices     HardwareDevices     `xml:"devices"`
}

type MemoryAllocation struct {
	Capacity int    `xml:",chardata"`
	Unit     string `xml:"unit,attr"`
}

type MemoryBacking struct {
	HugePages *struct{} `xml:"hugepages,omitempty"`
	Locked    *struct{} `xml:"locked,omitempty"`
}

type OperatingSystem struct {
	BootConfiguration BootConfiguration `xml:"type"`
}

type BootConfiguration struct {
	Architecture string `xml:"arch,attr"`
	MachineType  string `xml:"machine,attr"`
	KernelType   string `xml:",chardata"`
}

type SystemFeatures struct {
	AdvancedConfigurationPowerInterface xml.Name `xml:"acpi"`
	AdvancedProgrammableInterruptController xml.Name `xml:"apic"`
}

type ProcessorModel struct {
	ExecutionMode string       `xml:"mode,attr"`
	Topology      *CpuTopology `xml:"topology,omitempty"`
}

type CpuTopology struct {
	Sockets int `xml:"sockets,attr"`
	Cores   int `xml:"cores,attr"`
	Threads int `xml:"threads,attr"`
}

type SystemClock struct {
	TimeOffset string `xml:"offset,attr"`
}

type HardwareDevices struct {
	QemuEmulatorPath string            `xml:"emulator"`
	StorageDisks     []VirtualDisk     `xml:"disk"`
	NetworkInterface NetworkInterface  `xml:"interface"`
	GraphicsAdapter  GraphicsAdapter   `xml:"graphics"`
	VideoController  VideoController   `xml:"video"`
	MemBalloon       *MemBalloon       `xml:"memballoon,omitempty"`
}

type MemBalloon struct {
	Model string `xml:"model,attr"`
}

type VirtualDisk struct {
	StorageType   string          `xml:"type,attr"`
	DeviceType    string          `xml:"device,attr"`
	DiskDriver    DiskDriver      `xml:"driver"`
	SourceFile    DiskSourceFile  `xml:"source"`
	TargetDevice  DiskTargetDevice `xml:"target"`
}

type DiskDriver struct {
	DriverName string `xml:"name,attr"`
	FormatType string `xml:"type,attr"`
}

type DiskSourceFile struct {
	FilePath string `xml:"file,attr"`
}

type DiskTargetDevice struct {
	DevicePrefix string `xml:"dev,attr"`
	BusType      string `xml:"bus,attr"`
}

type NetworkInterface struct {
	InterfaceType string          `xml:"type,attr"`
	NetworkSource NetworkSource   `xml:"source"`
	HardwareModel InterfaceModel  `xml:"model"`
}

type NetworkSource struct {
	NetworkName string `xml:"network,attr"`
}

type InterfaceModel struct {
	ModelType string `xml:"type,attr"`
}

type GraphicsAdapter struct {
	ProtocolType   string        `xml:"type,attr"`
	AutoPortConfig string        `xml:"autoport,attr"`
	ListenAddress  ListenAddress `xml:"listen"`
}

type ListenAddress struct {
	AddressType string `xml:"type,attr"`
}

type VideoController struct {
	VideoModel VideoModel `xml:"model"`
}

type VideoModel struct {
	ModelType string `xml:"type,attr"`
}

func GenerateVirtualMachineXmlDefinition(
	name string,
	ramMegabytes int,
	cpuCores int,
	diskImageFullPath string,
	isoImageFullPath string,
) (
	string, 
	error,
) {
	qemuBinaryPath := "/usr/bin/qemu-system-x86_64"
	if runtime.GOOS == "darwin" {
		qemuBinaryPath = "/usr/local/bin/qemu-system-x86_64"
	}

	config := VirtualMachineConfiguration{
		HypervisorType: "kvm",
		MachineName:    name,
		MemoryAllocation: MemoryAllocation{
			Capacity: ramMegabytes * 1024,
			Unit:     "KiB",
		},
		CurrentMemory: MemoryAllocation{
			Capacity: ramMegabytes * 1024,
			Unit:     "KiB",
		},
		VirtualProcessorCount: cpuCores,
		OperatingSystem: OperatingSystem{
			BootConfiguration: BootConfiguration{
				Architecture: "x86_64",
				MachineType:  "q35",
				KernelType:   "hvm",
			},
		},
		SystemFeatures: SystemFeatures{
			AdvancedConfigurationPowerInterface: xml.Name{Local: "acpi"},
			AdvancedProgrammableInterruptController: xml.Name{Local: "apic"},
		},
		ProcessorModel: ProcessorModel{
			ExecutionMode: "host-passthrough",
			Topology: &CpuTopology{
				Sockets: 1,
				Cores:   cpuCores,
				Threads: 1,
			},
		},
		SystemClock: SystemClock{TimeOffset: "utc"},
		HardwareDevices: HardwareDevices{
			QemuEmulatorPath: qemuBinaryPath,
			StorageDisks: []VirtualDisk{
				{
					StorageType: "file",
					DeviceType:  "disk",
					DiskDriver: DiskDriver{
						DriverName: "qemu",
						FormatType: "qcow2",
					},
					SourceFile: DiskSourceFile{FilePath: diskImageFullPath},
					TargetDevice: DiskTargetDevice{
						DevicePrefix: "vda",
						BusType:      "virtio",
					},
				},
				{
					StorageType: "file",
					DeviceType:  "cdrom",
					DiskDriver: DiskDriver{
						DriverName: "qemu",
						FormatType: "raw",
					},
					SourceFile: DiskSourceFile{FilePath: isoImageFullPath},
					TargetDevice: DiskTargetDevice{
						DevicePrefix: "sda",
						BusType:      "sata",
					},
				},
			},
			NetworkInterface: NetworkInterface{
				InterfaceType: "network",
				NetworkSource: NetworkSource{NetworkName: "default"},
				HardwareModel: InterfaceModel{ModelType: "virtio"},
			},
			GraphicsAdapter: GraphicsAdapter{
				ProtocolType:   "spice",
				AutoPortConfig: "yes",
				ListenAddress:  ListenAddress{AddressType: "address"},
			},
			VideoController: VideoController{
				VideoModel: VideoModel{ModelType: "virtio"},
			},
			MemBalloon: &MemBalloon{Model: "virtio"},
		},
	}

	xmlBytes, err := xml.MarshalIndent(config, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed_to_marshal_virtual_machine_xml: %w", err)
	}

	return string(xmlBytes), nil
}
