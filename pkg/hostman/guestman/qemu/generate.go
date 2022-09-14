// Copyright 2019 Yunion
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package qemu

import (
	"fmt"
	"strings"

	"yunion.io/x/pkg/errors"
	"yunion.io/x/pkg/utils"

	api "yunion.io/x/onecloud/pkg/apis/compute"
	"yunion.io/x/onecloud/pkg/hostman/guestman/desc"
)

type Monitor struct {
	Id   string
	Port uint
	Mode string
}

func optionsToString(options map[string]string) string {
	var cmd string
	for key, value := range options {
		if value != "" {
			cmd += fmt.Sprintf(",%s=%s", key, value)
		} else {
			cmd += fmt.Sprintf(",%s", key)
		}
	}
	return cmd
}

func generatePCIDeviceOption(dev *desc.PCIDevice) string {
	cmd := fmt.Sprintf(
		"-device %s,id=%s", dev.DevType, dev.Id,
	)
	if dev.PCIAddr != nil {
		cmd += fmt.Sprintf(",bus=%s,addr=%s", dev.BusStr(), dev.SlotFunc())
		if dev.Multi != nil {
			cmd += fmt.Sprintf(",%s", dev.MultiFunction())
		}
	}

	cmd += optionsToString(dev.Options)
	return cmd
}

func chardevOption(c *desc.CharDev) string {
	cmd := fmt.Sprintf("-chardev %s,id=%s", c.Backend, c.Id)
	if c.Name != "" {
		cmd += fmt.Sprintf(",name=%s", c.Name)
	}
	cmd += optionsToString(c.Options)
	return cmd
}

func virtSerialPortOption(p *desc.VirtSerialPort, bus string) string {
	cmd := fmt.Sprintf(
		"-device virtserialport,bus=%s.0,chardev=%s,name=%s",
		bus, p.Chardev, p.Name,
	)
	cmd += optionsToString(p.Options)
	return cmd
}

func usbControllerOption(usb *desc.UsbController) string {
	cmd := generatePCIDeviceOption(usb.PCIDevice)
	if usb.MasterBus != nil {
		cmd += fmt.Sprintf(
			",masterbus=%s.0,firstport=%d",
			usb.MasterBus.Masterbus, usb.MasterBus.Port,
		)
	}
	return cmd
}

func usbRedirOptions(usbredir *desc.UsbRedirctDesc) []string {
	opts := make([]string, 0)
	opts = append(opts, usbControllerOption(usbredir.EHCI1))
	opts = append(opts, usbControllerOption(usbredir.UHCI1))
	opts = append(opts, usbControllerOption(usbredir.UHCI2))
	opts = append(opts, usbControllerOption(usbredir.UHCI3))
	opts = append(opts, chardevOption(usbredir.UsbRedirDev1.Source))
	opts = append(opts, fmt.Sprintf("-device usb-redir,chardev=%s,id=%s", usbredir.UsbRedirDev1.Source.Id, usbredir.UsbRedirDev1.Id))
	opts = append(opts, chardevOption(usbredir.UsbRedirDev2.Source))
	opts = append(opts, fmt.Sprintf("-device usb-redir,chardev=%s,id=%s", usbredir.UsbRedirDev2.Source.Id, usbredir.UsbRedirDev2.Id))

	return opts
}

func generateSpiceOptions(port uint, spice *desc.SSpiceDesc) []string {
	opts := make([]string, 0)

	// spice
	spiceCmd := fmt.Sprintf("-spice port=%d", port)
	spiceCmd += optionsToString(spice.Options)
	opts = append(opts, spiceCmd)

	// intel-hda and codec hda-duplex
	opts = append(opts, generatePCIDeviceOption(spice.IntelHDA.PCIDevice))
	codec := spice.IntelHDA.Codec
	opts = append(opts,
		fmt.Sprintf("-device %s,id=%s,bus=%s.0,cad=%d",
			codec.Type, codec.Id, spice.IntelHDA.Id, codec.Cad),
	)

	// serial port
	opts = append(opts, generatePCIDeviceOption(spice.VdagentSerial.PCIDevice))
	opts = append(opts, chardevOption(spice.Vdagent))
	opts = append(opts, virtSerialPortOption(spice.VdagentSerialPort, spice.VdagentSerial.Id))

	// usb redirct
	opts = append(opts, usbRedirOptions(spice.UsbRedirct)...)
	return opts
}

func generatePciControllerOptions(controllers []*desc.PCIController) []string {
	opts := make([]string, 0)
	for i := 0; i < len(controllers); i++ {
		switch controllers[i].CType {
		case desc.CONTROLLER_TYPE_PCIE_TO_PCI_BRIDGE:
			opts = append(opts,
				fmt.Sprintf(
					"-device pcie-pci-bridge,id=pci.%d,bus=pcie.%d,addr=0x%02x",
					i, controllers[i].Bus, controllers[i].Slot,
				))
		case desc.CONTROLLER_TYPE_PCI_BRIDGE:
			opts = append(opts,
				fmt.Sprintf(
					"-device pci-bridge,id=pci.%d,bus=pci.%d,chassis_nr=%d,addr=0x%02x",
					i, controllers[i].Bus, i, controllers[i].Slot),
			)
		case desc.CONTROLLER_TYPE_PCIE_ROOT_PORT:
			opts = append(opts,
				fmt.Sprintf(
					"-device ioh3420,id=pci.%d,chassis=%d,bus=pcie.%d,addr=0x%02x",
					i, i, controllers[i].Bus, controllers[i].Slot,
				),
			)
		}
	}
	return opts
}

func generateMemoryOption(sizeMB int64, memDesc *desc.SGuestMem) string {
	return fmt.Sprintf(
		"-m %dM,slots=%d,maxmem=%dM",
		sizeMB, memDesc.Slots, memDesc.MaxMem,
	)
}

func generateMachineOption(machine string, machineDesc *desc.SGuestMachine) string {
	cmd := fmt.Sprintf("-machine %s,accel=%s", machine, machineDesc.Accel)
	if machineDesc.GicVersion != nil {
		cmd += fmt.Sprintf(",gic-version=%s", *machineDesc.GicVersion)
	}

	return cmd
}

func generateSMPOption(cpu *desc.SGuestCpu) string {
	return fmt.Sprintf(
		"-smp cpus=%d,sockets=%d,cores=%d,maxcpus=%d",
		cpu.Cpus, cpu.Sockets, cpu.Cores, cpu.MaxCpus,
	)
}

func generateCPUOption(cpu *desc.SGuestCpu) string {
	cmd := fmt.Sprintf("-cpu %s", cpu.Model)
	for feat, enable := range cpu.Features {
		if enable {
			cmd += "," + feat + "=on"
		} else {
			cmd += "," + feat + "=off"
		}
	}
	if len(cpu.Vendor) > 0 {
		cmd += fmt.Sprintf(",vendor=%s", cpu.Vendor)
	}
	if len(cpu.Level) > 0 {
		cmd += fmt.Sprintf(",level=%s", cpu.Level)
	}
	return cmd
}

func getMonitorOptions(drvOpt QemuOptions, input *Monitor) []string {
	if input == nil {
		return nil
	}
	idDev := input.Id + "dev"
	opts := []string{
		drvOpt.MonitorChardev(idDev, input.Port, "127.0.0.1"),
		drvOpt.Mon(idDev, input.Id, input.Mode),
	}
	return opts
}

func generateDisksOptions(drvOpt QemuOptions, disks []*desc.SGuestDisk, isEncrypt bool) []string {
	opts := make([]string, 0)
	for _, disk := range disks {
		opts = append(opts,
			getDiskDriveOption(drvOpt, disk, isEncrypt),
			getDiskDeviceOption(drvOpt, disk),
		)
	}
	return opts
}

func getDiskDriveOption(drvOpt QemuOptions, disk *desc.SGuestDisk, isEncrypt bool) string {
	format := disk.Format
	diskIndex := disk.Index
	cacheMode := disk.CacheMode
	aioMode := disk.AioMode

	opt := fmt.Sprintf("file=$DISK_%d", diskIndex)
	opt += ",if=none"
	opt += fmt.Sprintf(",id=drive_%d", diskIndex)
	if len(format) == 0 || format == "qcow2" {
		// pass    # qemu will automatically detect image format
	} else if format == "raw" {
		opt += ",format=raw"
	}
	opt += fmt.Sprintf(",cache=%s", cacheMode)
	if isLocalStorage(disk) {
		opt += fmt.Sprintf(",aio=%s", aioMode)
	}
	if len(disk.Url) > 0 { // # a remote file backed image
		opt += ",copy-on-read=on"
	}
	if isLocalStorage(disk) {
		opt += ",file.locking=off"
	}
	if isEncrypt {
		opt += ",encrypt.format=luks,encrypt.key-secret=sec0"
	}
	// #opt += ",media=disk"
	return drvOpt.Drive(opt)
}

func isLocalStorage(disk *desc.SGuestDisk) bool {
	if disk.StorageType == api.STORAGE_LOCAL || len(disk.StorageType) == 0 {
		return true
	} else {
		return false
	}
}

func getDiskDeviceOption(optDrv QemuOptions, disk *desc.SGuestDisk) string {
	diskIndex := disk.Index
	diskDriver := disk.Driver
	numQueues := disk.NumQueues
	isSsd := disk.IsSSD

	if numQueues == 0 {
		numQueues = 4
	}

	var opt = ""
	opt += GetDiskDeviceModel(diskDriver)
	opt += fmt.Sprintf(",drive=drive_%d", diskIndex)
	if diskDriver == DISK_DRIVER_VIRTIO {
		// virtio-blk
		if disk.Pci != nil {
			opt += fmt.Sprintf(",bus=%s,addr=0x%s", disk.Pci.BusStr(), disk.Pci.SlotFunc())
		}
		// opt += fmt.Sprintf(",num-queues=%d,vectors=%d,iothread=iothread0", numQueues, numQueues+1)
		opt += ",iothread=iothread0"
	} else if utils.IsInStringArray(diskDriver, []string{DISK_DRIVER_SCSI, DISK_DRIVER_PVSCSI}) {
		opt += ",bus=scsi.0"
	} else if diskDriver == DISK_DRIVER_IDE {
		opt += fmt.Sprintf(",bus=ide.%d,unit=%d", diskIndex/2, diskIndex%2)
	} else if diskDriver == DISK_DRIVER_SATA {
		opt += fmt.Sprintf(",bus=ide.%d", diskIndex)
	}
	opt += fmt.Sprintf(",id=drive_%d", diskIndex)
	if isSsd {
		opt += ",rotation_rate=1"
	}
	return optDrv.Device(opt)
}

func generateCdromOptions(optDrv QemuOptions, cdrom *desc.SGuestCdrom) []string {
	opts := make([]string, 0)

	//cdromDriveId := cdrom
	driveOpt := fmt.Sprintf("id=%s", cdrom.Id)
	driveOpt += optionsToString(cdrom.DriveOptions)

	var cdromPath = cdrom.Path
	if len(cdromPath) > 0 {
		driveOpt += fmt.Sprintf(",file=%s", cdromPath)
	}

	if cdrom.Ide != nil {
		opts = append(opts, optDrv.Drive(driveOpt))
		devOpt := fmt.Sprintf("%s,drive=%s,bus=ide.1",
			cdrom.Ide.DevType, cdrom.Id)
		// TODO: ,bus=ide.%d,unit=%d
		//, cdrom.Ide.Bus, cdrom.Ide.Unit)
		opts = append(opts, optDrv.Device(devOpt))
	} else if cdrom.Scsi != nil {
		if len(cdromPath) > 0 {
			opts = append(opts, optDrv.Drive(driveOpt))

			devOpt := fmt.Sprintf("%s,drive=%s", cdrom.Scsi.DevType, cdrom.Id)
			devOpt += optionsToString(cdrom.Scsi.Options)

			opts = append(opts, optDrv.Device(devOpt))
		}
	}
	return opts
}

func GetDiskDeviceModel(driver string) string {
	if driver == DISK_DRIVER_VIRTIO {
		return "virtio-blk-pci"
	} else if utils.IsInStringArray(driver, []string{DISK_DRIVER_SCSI, DISK_DRIVER_PVSCSI}) {
		return "scsi-hd"
	} else if driver == DISK_DRIVER_IDE {
		return "ide-hd"
	} else if driver == DISK_DRIVER_SATA {
		return "ide-drive"
	} else {
		return "None"
	}
}

func generateNicOptions(drvOpt QemuOptions, input *GenerateStartOptionsInput) ([]string, error) {
	opts := make([]string, 0)
	nics := input.GuestDesc.Nics

	for idx := range nics {
		netDevOpt, err := getNicNetdevOption(drvOpt, nics[idx], input.IsKVMSupport)
		if err != nil {
			return nil, errors.Wrapf(err, "getNicNetdevOption %v", nics[idx])
		}
		opts = append(opts,
			netDevOpt,
			// aarch64 with addr lead to:
			// virtio_net: probe of virtioN failed with error -22
			getNicDeviceOption(drvOpt, nics[idx], input))
	}
	return opts, nil
}

func getNicNetdevOption(drvOpt QemuOptions, nic *desc.SGuestNetwork, isKVMSupport bool) (string, error) {
	if nic.Ifname == "" {
		return "", errors.Error("ifname is empty")
	}
	if nic.UpscriptPath == "" {
		return "", errors.Error("upscript_path is empty")
	}
	if nic.DownscriptPath == "" {
		return "", errors.Error("downscript_path is empty")
	}

	opt := "-netdev type=tap"
	opt += fmt.Sprintf(",id=%s", nic.Ifname)
	opt += fmt.Sprintf(",ifname=%s", nic.Ifname)
	if nic.Driver == "virtio" && isKVMSupport {
		opt += ",vhost=on,vhostforce=off"
		if nic.NumQueues > 1 {
			opt += fmt.Sprintf(",queues=%d", nic.NumQueues)
		}
	}
	opt += fmt.Sprintf(",script=%s", nic.UpscriptPath)
	opt += fmt.Sprintf(",downscript=%s", nic.DownscriptPath)
	return opt, nil
}

func getNicDeviceOption(
	drvOpt QemuOptions,
	nic *desc.SGuestNetwork,
	input *GenerateStartOptionsInput,
) string {
	cmd := generatePCIDeviceOption(nic.Pci)
	cmd += fmt.Sprintf(",netdev=%s", nic.Ifname)
	cmd += fmt.Sprintf(",mac=%s", nic.Mac)

	if nic.Driver == "virtio" {
		if nic.NumQueues > 1 {
			cmd += fmt.Sprintf(",mq=on")
		}
		if nic.Vectors != nil {
			cmd += fmt.Sprintf(",vectors=%d", *nic.Vectors)
		}
		cmd += fmt.Sprintf("$(nic_speed %d)", nic.Bw)
		if nic.Bridge == input.OVNIntegrationBridge {
			cmd += fmt.Sprintf("$(nic_mtu %q)", nic.Bridge)
		}
	}
	return cmd
}

func GetNicDeviceModel(name string) string {
	if name == "virtio" {
		return "virtio-net-pci"
	} else if name == "e1000" {
		return "e1000-82545em"
	} else {
		return name
	}
}

func generateUsbDeviceOption(usbControllerId string, usb *desc.UsbDevice) string {
	cmd := fmt.Sprintf("-device %s,bus=%s.0", usb.DevType, usbControllerId)
	cmd += optionsToString(usb.Options)
	return cmd
}

func generateVfioDeviceOption(devs []*desc.VFIODevice) []string {
	opts := make([]string, 0)

	for i := 0; i < len(devs); i++ {
		cmd := generatePCIDeviceOption(devs[i].PCIDevice)
		cmd += fmt.Sprintf(",host=%s", devs[i].HostAddr)
		if devs[i].XVga {
			cmd += ",x-vga=on"
		}
		opts = append(opts, cmd)
	}
	return opts
}

func generateIsolatedDeviceOptions(guestDesc *desc.SGuestDesc) []string {
	opts := make([]string, 0)
	for i := 0; i < len(guestDesc.IsolatedDevices); i++ {
		if guestDesc.IsolatedDevices[i].Usb != nil {
			opts = append(opts,
				generateUsbDeviceOption(guestDesc.Usb.Id, guestDesc.IsolatedDevices[i].Usb),
			)
		} else if len(guestDesc.IsolatedDevices[i].VfioDevs) > 0 {
			opts = append(opts,
				generateVfioDeviceOption(guestDesc.IsolatedDevices[i].VfioDevs)...,
			)
		}
	}
	return opts
}

func generateObjectOption(o *desc.Object) string {
	cmd := fmt.Sprintf("-object %s,id=%s", o.ObjType, o.Id)
	cmd += optionsToString(o.Options)
	return cmd
}

func getRNGRandomOptions(rng *desc.SGuestRng) []string {
	cmd := generatePCIDeviceOption(rng.PCIDevice)
	cmd += fmt.Sprintf(",rng=%s", rng.RngRandom.Id)

	return []string{
		generateObjectOption(rng.RngRandom),
		cmd,
	}
}

func generateQgaOptions(guestDesc *desc.SGuestDesc) []string {
	opts := make([]string, 0)
	opts = append(opts, chardevOption(guestDesc.Qga.Socket))
	opts = append(opts, virtSerialPortOption(guestDesc.Qga.SerialPort, guestDesc.VirtioSerial.Id))
	return opts
}

func generateISASerialOptions(isaSerial *desc.SGuestIsaSerial) []string {
	opts := make([]string, 0)
	opts = append(opts, chardevOption(isaSerial.Pty))
	opts = append(opts, fmt.Sprintf("-device isa-serial,chardev=%s,id=%s", isaSerial.Pty.Id, isaSerial.Id))
	return opts
}

func generatePvpanicDeviceOption(pvpanic *desc.SGuestPvpanic) string {
	return fmt.Sprintf("-device pvpanic,id=%s,ioport=0x%x", pvpanic.Id, pvpanic.Ioport)
}

func getMigrateOptions(drvOpt QemuOptions, input *GenerateStartOptionsInput) []string {
	opts := make([]string, 0)
	if input.NeedMigrate {
		if input.LiveMigrateUseTLS {
			opts = append(opts, fmt.Sprintf("-incoming defer"))
		} else {
			opts = append(opts, fmt.Sprintf("-incoming tcp:0:%d", input.LiveMigratePort))
		}
	} else if input.GuestDesc.IsSlave {
		opts = append(opts, fmt.Sprintf("-incoming tcp:0:%d", input.LiveMigratePort))
	}
	return opts
}

type GenerateStartOptionsInput struct {
	QemuVersion Version
	QemuArch    Arch

	GuestDesc    *desc.SGuestDesc
	IsKVMSupport bool

	EnableUUID       bool
	OsName           string
	HugepagesEnabled bool
	EnableMemfd      bool

	OVNIntegrationBridge string
	Devices              []string
	OVMFPath             string
	VNCPort              uint
	VNCPassword          bool
	EnableLog            bool
	LogPath              string
	HMPMonitor           *Monitor
	QMPMonitor           *Monitor
	IsVdiSpice           bool
	SpicePort            uint
	PidFilePath          string
	HomeDir              string
	ExtraOptions         []string
	EnableRNGRandom      bool
	EnableSerialDevice   bool
	NeedMigrate          bool
	LiveMigratePort      uint
	LiveMigrateUseTLS    bool
	EnablePvpanic        bool

	EncryptKeyPath string
}

func GenerateStartOptions(
	input *GenerateStartOptionsInput,
) (string, error) {
	drv, ok := GetCommand(input.QemuVersion, input.QemuArch)
	if !ok {
		return "", errors.Errorf("Qemu comand driver %s %s not registered", input.QemuVersion, input.QemuArch)
	}
	drvOpt := drv.GetOptions()

	opts := make([]string, 0)

	// generate cpu options
	cpuOpt := generateCPUOption(input.GuestDesc.CpuDesc)
	opts = append(opts, drvOpt.FreezeCPU(), cpuOpt)

	if input.EnableLog {
		opts = append(opts, drvOpt.Log(input.EnableLog, input.LogPath))
	}

	// TODO hmp - -
	opts = append(opts, getMonitorOptions(drvOpt, input.HMPMonitor)...)
	if input.QMPMonitor != nil {
		opts = append(opts, getMonitorOptions(drvOpt, input.QMPMonitor)...)
	}

	opts = append(opts,
		drvOpt.RTC(),
		drvOpt.Daemonize(),
		drvOpt.Nodefaults(),
		drvOpt.Nodefconfig(),
		// drvOpt.NoKVMPitReinjection(),
		drvOpt.Global(),
		generateMachineOption(input.GuestDesc.Machine, input.GuestDesc.MachineDesc),
		drvOpt.KeyboardLayoutLanguage("en-us"),
		generateSMPOption(input.GuestDesc.CpuDesc),
		drvOpt.Name(input.GuestDesc.Name),
		drvOpt.UUID(input.EnableUUID, input.GuestDesc.Uuid),
		generateMemoryOption(input.GuestDesc.Mem, input.GuestDesc.MemDesc),
	)

	var memDev string
	if input.HugepagesEnabled {
		memDev = drvOpt.MemPath(
			uint64(input.GuestDesc.Mem),
			fmt.Sprintf("/dev/hugepages/%s", input.GuestDesc.Uuid),
		)
	} else if input.EnableMemfd {
		memDev = drvOpt.MemFd(uint64(input.GuestDesc.Mem))
	} else {
		memDev = drvOpt.MemDev(uint64(input.GuestDesc.Mem))
	}
	opts = append(opts, memDev)

	// bootOrder
	enableMenu := false
	if input.GuestDesc.Cdrom != nil && input.GuestDesc.Cdrom.Path != "" {
		enableMenu = true
	}
	opts = append(opts, drvOpt.Boot(input.GuestDesc.BootOrder, enableMenu))

	// bios
	if input.GuestDesc.Bios == BIOS_UEFI {
		if input.OVMFPath == "" {
			return "", errors.Errorf("input OVMF path is empty")
		}
		fmOpt, err := drvOpt.BIOS(input.OVMFPath, input.HomeDir)
		if err != nil {
			return "", errors.Wrap(err, "bios option")
		}
		opts = append(opts, fmOpt)
	}

	if input.OsName == OS_NAME_MACOS {
		opts = append(opts, drvOpt.Device("isa-applesmc,osk=ourhardworkbythesewordsguardedpleasedontsteal(c)AppleComputerInc"))
	}

	if input.GuestDesc.Vga != "none" {
		opts = append(opts, generatePCIDeviceOption(input.GuestDesc.VgaDevice.PCIDevice))
	}

	// vdi spice
	if input.IsVdiSpice {
		opts = append(opts, generateSpiceOptions(input.SpicePort, input.GuestDesc.VdiDevice.Spice)...)
	} else {
		opts = append(opts, drvOpt.VNC(input.VNCPort, input.VNCPassword))
	}

	// iothread object
	opts = append(opts, drvOpt.Object("iothread", map[string]string{"id": "iothread0"}))

	isEncrypt := false
	if len(input.EncryptKeyPath) > 0 {
		opts = append(opts, drvOpt.Object("secret", map[string]string{"id": "sec0", "file": input.EncryptKeyPath, "format": "base64"}))
		isEncrypt = true
	}

	if input.GuestDesc.VirtioSerial != nil {
		opts = append(opts, generatePCIDeviceOption(input.GuestDesc.VirtioSerial.PCIDevice))
	}

	opts = append(opts, generatePciControllerOptions(input.GuestDesc.PCIControllers)...)

	if input.GuestDesc.VirtioScsi != nil {
		opts = append(opts, generatePCIDeviceOption(input.GuestDesc.VirtioScsi.PCIDevice))
	} else if input.GuestDesc.PvScsi != nil {
		opts = append(opts, generatePCIDeviceOption(input.GuestDesc.PvScsi.PCIDevice))
	}
	// generate disk options
	opts = append(opts, generateDisksOptions(drvOpt, input.GuestDesc.Disks, isEncrypt)...)

	// cdrom
	opts = append(opts, generateCdromOptions(drvOpt, input.GuestDesc.Cdrom)...)

	// generate nics
	nicOpts, err := generateNicOptions(drvOpt, input)
	if err != nil {
		return "", errors.Wrap(err, "generateNicOptions")
	}
	opts = append(opts, nicOpts...)

	// isolated devices
	// USB 3.0
	opts = append(opts, generatePCIDeviceOption(input.GuestDesc.Usb.PCIDevice))
	// enable machine USB emulation
	opts = append(opts, drvOpt.USB())
	for _, device := range input.Devices {
		opts = append(opts, drvOpt.Device(device))
	}
	if len(input.GuestDesc.IsolatedDevices) > 0 {
		opts = append(opts, generateIsolatedDeviceOptions(input.GuestDesc)...)
	}

	// pidfile
	opts = append(opts, drvOpt.Pidfile(input.PidFilePath))

	// extra options
	if len(input.ExtraOptions) != 0 {
		opts = append(opts, input.ExtraOptions...)
	}

	// qga
	// opts = append(opts, drvOpt.QGA(input.HomeDir)...)
	if input.GuestDesc.Qga != nil {
		opts = append(opts, generateQgaOptions(input.GuestDesc)...)
	}

	// random device
	if input.GuestDesc.Rng != nil {
		opts = append(opts, getRNGRandomOptions(input.GuestDesc.Rng)...)
	}

	// serial device
	if input.GuestDesc.IsaSerial != nil {
		opts = append(opts, generateISASerialOptions(input.GuestDesc.IsaSerial)...)
	}

	// migrate options
	opts = append(opts, getMigrateOptions(drvOpt, input)...)

	// pvpanic device
	if input.GuestDesc.Pvpanic != nil {
		opts = append(opts, generatePvpanicDeviceOption(input.GuestDesc.Pvpanic))
	}

	return strings.Join(opts, " "), nil
}
