package ova

import (
	"encoding/xml"
	"fmt"
	"strconv"

	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	cnv "kubevirt.io/api/core/v1"
)

// ResourceType
const (
	OtherResourceType                 = "1"
	ComputerSystemResourceType        = "2"
	ProcessorResourceType             = "3"
	MemoryResoureceType               = "4"
	IDEControllerResourceType         = "5"
	ParallelSCSIHBAResourceType       = "6"
	FCHBAResourceType                 = "7"
	ISCSIHBAResourceType              = "8"
	IBHCAResourceType                 = "9"
	EthernetAdapterResourceType       = "10"
	OtherNetworkAdapterResourceType   = "11"
	IOSlotResourceType                = "12"
	IODeviceResourceType              = "13"
	FloppyDriveResourceType           = "14"
	CDDriveResourceType               = "15"
	DVDDriveResourceType              = "16"
	DiskDriveResourceType             = "17"
	TapeDriveResourceType             = "18"
	StorageExentResourceType          = "19"
	OtherStorageDeviceResourceType    = "20"
	SerialPortResourceType            = "21"
	ParallelPortResourceType          = "22"
	USBControllerResourceType         = "23"
	GraphicsControllerResourceType    = "24"
	IEEE1394ControllerResourceType    = "25"
	PartitionableUnitResourceType     = "26"
	BasePartitionableUnitResourceType = "27"
	PowerSupplyResourceType           = "28"
	CoolingDeviceResourceType         = "29"
	PCIControllerResourceType         = "30"
	PS2ControllerResourceType         = "31"
	SIOControllerResourceType         = "32"
	KeyboardResourceType              = "33"
	PointingDeviceResourceType        = "34"
)

// Map of osinfo ids to vmware guest ids.
var osMap = map[string]string{
	"centos5.11":            "centos64Guest",
	"centos6.10":            "centos6_64Guest",
	"centos7.0":             "centos7_64Guest",
	"centos8":               "centos8_64Guest",
	"debian4":               "debian4_64Guest",
	"debian5":               "debian5_64Guest",
	"debian6":               "debian6_64Guest",
	"debian7":               "debian7_64Guest",
	"debian8":               "debian8_64Guest",
	"debian9":               "debian9_64Guest",
	"debian10":              "debian10_64Guest",
	"fedora":                "fedora64Guest",
	"fedora31":              "fedora64Guest",
	"linux":                 "genericLinuxGuest",
	"rhel6.10":              "rhel6_64Guest",
	"rhel7.7":               "rhel7_64Guest",
	"rhel8.1":               "rhel8_64Guest",
	"ubuntu18.04":           "ubuntu64Guest",
	"win2k":                 "win2000ServGuest",
	"win7":                  "windows7Guest",
	"win2k8r2":              "windows7Server64Guest",
	"win8":                  "windows8_64Guest",
	"win2k12r2":             "windows8Server64Guest",
	"win10":                 "windows9_64Guest",
	"win2k19":               "windows9Server64Guest",
	"windows9Server64Guest": "win2k19",
}

// xml struct
type Item struct {
	AllocationUnits     string          `xml:"rasd:AllocationUnits,omitempty"`
	Description         string          `xml:"rasd:Description,omitempty"`
	ElementName         string          `xml:"rasd:ElementName"`
	InstanceID          string          `xml:"rasd:InstanceID"`
	ResourceType        string          `xml:"rasd:ResourceType"`
	VirtualQuantity     int32           `xml:"rasd:VirtualQuantity"`
	Address             string          `xml:"rasd:Address,omitempty"`
	ResourceSubType     string          `xml:"rasd:ResourceSubType,omitempty"`
	Parent              string          `xml:"rasd:Parent,omitempty"`
	HostResource        string          `xml:"rasd:HostResource,omitempty"`
	Connection          string          `xml:"rasd:Connection,omitempty"`
	Configs             []VirtualConfig `xml:"vmw:Config"`
	CoresPerSocket      string          `xml:"CoresPerSocket"`
	Required            string          `xml:"ovf:required,attr,omitempty"`
	AutomaticAllocation bool            `xml:"rasd:AutomaticAllocation,omitempty"`
}

type VirtualConfig struct {
	XMLName  xml.Name `xml:"http://www.vmware.com/schema/ovf Config"`
	Required string   `xml:"ovf:required,attr"`
	Key      string   `xml:"vmw:key,attr"`
	Value    string   `xml:"vmw:value,attr"`
}

type VirtualHardwareSection struct {
	Info    string              `xml:"Info"`
	Items   []Item              `xml:"Item"`
	Configs []VirtualConfig     `xml:"Config"`
	System  VirtualSystemConfig `xml:"System"`
}

type VirtualSystemConfig struct {
	ElementName             string `xml:"vssd:ElementName"`
	InstanceID              string `xml:"vssd:InstanceID"`
	VirtualSystemIdentifier string `xml:"vssd:VirtualSystemIdentifier"`
}

type References struct {
	File []File `xml:"File"`
}

type File struct {
	ID   string `xml:"ovf:id,attr"`
	Href string `xml:"ovf:href,attr"`
	Size string `xml:"ovf:size,attr"`
}

type DiskSection struct {
	XMLName xml.Name `xml:"DiskSection"`
	Info    string   `xml:"Info"`
	Disks   []Disk   `xml:"Disk"`
}

type Disk struct {
	XMLName                 xml.Name `xml:"Disk"`
	Capacity                int64    `xml:"ovf:capacity,attr"`
	CapacityAllocationUnits string   `xml:"ovf:capacityAllocationUnits,attr"`
	DiskId                  string   `xml:"ovf:diskId,attr"`
	FileRef                 string   `xml:"ovf:fileRef,attr"`
	Format                  string   `xml:"ovf:format,attr"`
	PopulatedSize           int64    `xml:"ovf:populatedSize,attr"`
}

type NetworkSection struct {
	XMLName  xml.Name  `xml:"NetworkSection"`
	Info     string    `xml:"Info"`
	Networks []Network `xml:"Network"`
}

type Network struct {
	XMLName     xml.Name `xml:"Network"`
	Name        string   `xml:"name,attr"`
	Description string   `xml:"Description"`
}

type VirtualSystem struct {
	ID                     string                 `xml:"ovf:id,attr"`
	Name                   string                 `xml:"Name"`
	OperatingSystemSection OperatingSystemSection `xml:"OperatingSystemSection"`
	HardwareSection        VirtualHardwareSection `xml:"VirtualHardwareSection"`
}

type OperatingSystemSection struct {
	Info        string `xml:"Info"`
	Description string `xml:"Description"`
	OsType      string `xml:"vmw:osType,attr"`
}

type Envelope struct {
	XMLName        xml.Name       `xml:"Envelope"`
	XMLAttrOvf     string         `xml:"xmlns:ovf,attr"`
	XMLAttrRasd    string         `xml:"xmlns:rasd,attr"`
	XMLAttrVssd    string         `xml:"xmlns:vssd,attr"`
	XMLAttrXsi     string         `xml:"xmlns:xsi,attr"`
	VirtualSystem  VirtualSystem  `xml:"VirtualSystem"`
	DiskSection    DiskSection    `xml:"DiskSection"`
	NetworkSection NetworkSection `xml:"NetworkSection"`
	References     References     `xml:"References"`
}

func writeOvf(vm *cnv.VirtualMachine) {
	env := &Envelope{
		XMLAttrOvf:     "http://schemas.dmtf.org/ovf/envelope/1",
		XMLAttrRasd:    "http://schemas.dmtf.org/wbem/wscim/1/cim-schema/2/CIM_ResourceAllocationSettingData",
		XMLAttrVssd:    "http://schemas.dmtf.org/wbem/wscim/1/cim-schema/2/CIM_VirtualSystemSettingData",
		XMLAttrXsi:     "http://www.w3.org/2001/XMLSchema-instance",
		VirtualSystem:  VirtualSystem{},
		DiskSection:    DiskSection{},
		NetworkSection: NetworkSection{},
		References:     References{},
	}

	spec := vm.Spec.Template.Spec
	instanceId := 0
	env.VirtualSystem = VirtualSystem{
		ID:   vm.Name,
		Name: vm.Name,
		OperatingSystemSection: OperatingSystemSection{
			Info:   "The operating system installed",
			OsType: osMap[vm.Spec.Template.ObjectMeta.Annotations["vm.kubevirt.io/os"]],
		},
	}
	cpuItem := writeCpu(spec.Domain.CPU, &instanceId)
	env.VirtualSystem.HardwareSection.Info = "Virtual hardware requirements"
	env.VirtualSystem.HardwareSection.System = VirtualSystemConfig{
		InstanceID:              strconv.Itoa(instanceId),
		ElementName:             "Virtual Hardware Family",
		VirtualSystemIdentifier: vm.Name,
	}
	instanceId += 1
	env.VirtualSystem.HardwareSection.Items = append(env.VirtualSystem.HardwareSection.Items, *cpuItem)
	env.VirtualSystem.HardwareSection.Configs = vmConfigProperties(spec)

	out, _ := xml.MarshalIndent(env, " ", " ")
	fmt.Println(xml.Header + string(out))
}

func writeCpu(cpu *cnv.CPU, instanceId *int) (item *Item) {
	numOfCpus := cpu.Cores * cpu.Sockets * cpu.Threads
	item = &Item{
		InstanceID:      strconv.Itoa(*instanceId),
		ResourceType:    ProcessorResourceType,
		Description:     "Number of Virtual CPUs",
		AllocationUnits: "hertz * 10^6",
		VirtualQuantity: int32(numOfCpus),
		ElementName:     fmt.Sprintf(strconv.Itoa(int(numOfCpus)) + " virtual CPU(s)"),
	}
	*instanceId += 1
	return
}

func writeMemory(vmSpec *cnv.VirtualMachineInstanceSpec, memory *cnv.Memory, instanceId *int) (item *Item, err error) {
	var mem resource.Quantity
	if memory.Guest != nil {
		mem = *memory.Guest
	} else {
		mem = vmSpec.Domain.Resources.Requests[core.ResourceMemory]
	}
	parsedMem, ok := mem.AsInt64()
	if !ok {
		err = fmt.Errorf("Couldn't parse the quantity of the memory")
		return
	}
	item = &Item{
		InstanceID:      strconv.Itoa(*instanceId),
		ResourceType:    MemoryResoureceType,
		Description:     "Memory Size",
		AllocationUnits: "byte * 2^20",
		VirtualQuantity: int32(parsedMem),
		ElementName:     fmt.Sprintf(strconv.Itoa(int(parsedMem)) + "MB of memory"),
	}
	*instanceId += 1
	return
}

func writeSATA(instanceId *int) (item *Item) {
	item = &Item{
		InstanceID:      strconv.Itoa(*instanceId),
		ResourceType:    OtherStorageDeviceResourceType,
		Description:     "SATA Controller",
		ResourceSubType: "vmware.sata.ahci",
		ElementName:     "SATA Controller 0",
	}
	*instanceId += 1
	return
}

func writeSCSI(instanceId *int) (item *Item) {
	item = &Item{
		InstanceID:      strconv.Itoa(*instanceId),
		ResourceType:    ParallelSCSIHBAResourceType,
		Description:     "SCSI Controller",
		ResourceSubType: "VirtualSCSI",
		ElementName:     "SCSI Controller 0",
	}
	*instanceId += 1
	return
}

func writeIDE(instanceId *int) (item *Item) {
	item = &Item{
		InstanceID:   strconv.Itoa(*instanceId),
		ResourceType: IDEControllerResourceType,
		Description:  "IDE Controller",
		ElementName:  "IDE 0",
	}
	*instanceId += 1
	return
}

func writeCD(parent, instanceId *int) (item *Item) {
	item = &Item{
		AutomaticAllocation: false,
		ElementName:         "CD/DVD Drive 1",
		InstanceID:          strconv.Itoa(*instanceId),
		Parent:              strconv.Itoa(*parent),
		ResourceSubType:     "vmware.cdrom.remotepassthrough", // can be also vmware.cdrom.iso
		ResourceType:        CDDriveResourceType,
	}
	*instanceId += 1
	return
}

func writeVideo(instanceId *int) (item *Item) {
	item = &Item{
		Required:            "false",
		AutomaticAllocation: false,
		InstanceID:          strconv.Itoa(*instanceId),
		ResourceType:        GraphicsControllerResourceType,
		ElementName:         "Video card",
		Configs: []VirtualConfig{ //since they are not required - maybe drop it
			{
				Required: "false",
				Key:      "enable3DSupport",
				Value:    "false",
			},
			{
				Required: "false",
				Key:      "use3dRenderer",
				Value:    "automatic",
			},
			{
				Required: "false",
				Key:      "useAutoDetect",
				Value:    "false",
			},
			{
				Required: "false",
				Key:      "videoRamSizeInKB",
				Value:    "8192", // or half?
			},
		},
	}
	*instanceId += 1
	return
}

func writeHardDrive(disk string, parentId, slot, instanceId *int) (item *Item) {
	item = &Item{
		InstanceID:   strconv.Itoa(*instanceId),
		ElementName:  fmt.Sprintf("Hard disk %d", *slot+1),
		Address:      strconv.Itoa(*slot),
		ResourceType: DiskDriveResourceType,
		Parent:       strconv.Itoa(*parentId),           // scsi controller
		HostResource: fmt.Sprintf("ovf:/disk/%s", disk), // disk == vmdisk1
	}
	*instanceId += 1
	return
}

func writeNetwork(mac string, nicNum int, instanceId *int) (item *Item) {
	item = &Item{
		AutomaticAllocation: true,
		Connection:          "VM Network",
		Description:         "VmxNet3 ethernet adapter on &quot;VM Network&quot;",
		ElementName:         fmt.Sprintf("Network adapter %d", nicNum+1),
		InstanceID:          strconv.Itoa(*instanceId),
		ResourceSubType:     "VmxNet3",
		ResourceType:        EthernetAdapterResourceType,
	}
	if mac != "" {
		item.Address = mac
	}
	*instanceId += 1
	return
}

func vmConfigProperties(spec cnv.VirtualMachineInstanceSpec) (configs []VirtualConfig) {
	bootloader := spec.Domain.Firmware.Bootloader
	firmware := "bios"
	if bootloader.EFI != nil {
		firmware = "efi"
	}
	configs = []VirtualConfig{
		{
			Required: "false",
			Key:      "firmware",
			Value:    firmware,
		},
		{
			Required: "false",
			Key:      "uuid",
			Value:    string(spec.Domain.Firmware.UUID),
		},
	}
	return
}

func networkSection(spec cnv.VirtualMachineInstanceSpec) (networkSection NetworkSection) {
	networkSection = NetworkSection{
		Info:     "The list of logical networks",
		Networks: []Network{},
	}

	for _, network := range spec.Networks {
		networkSection.Networks = append(networkSection.Networks, Network{Name: network.Name, Description: "The VM Network network"})
	}
	return
}

func diskSection(vmName string, spec cnv.VirtualMachineInstanceSpec) (diskSection DiskSection, references *References) {
	diskSection = DiskSection{
		Info:  "List of the virtual disks",
		Disks: []Disk{},
	}
	references = &References{File: []File{}}
	diskNum := 1
	for _, disk := range spec.Volumes {
		diskSection.Disks = append(diskSection.Disks, Disk{
			CapacityAllocationUnits: "byte",
			Format:                  "http://www.vmware.com/interfaces/specifications/vmdk.html#streamOptimized",
			DiskId:                  fmt.Sprintf("vmdisk%d", diskNum),
			Capacity:                2,//pvc req size
			FileRef:                 fmt.Sprintf("file%d", diskNum),
			//PopulatedSize: ,NEED TO MEASURE
		})
		references.File = append(references.File, File{
			ID:   fmt.Sprintf("file%d", diskNum),
			Href: fmt.Sprintf("%s-disk%d.vmdk", vmName, diskNum),
		})
	}
}
