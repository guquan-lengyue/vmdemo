// Package kvm virsh 命令调用工具
package kvm

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

// VMInfo 表示虚拟机的基本信息
type VMInfo struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	State string `json:"state"`
	Vnc   string `json:"vnc"`
}

// ExecVirshCommand 执行 virsh 命令并返回输出结果
func ExecVirshCommand(args ...string) (string, error) {
	cmd := exec.Command("virsh", args...)
	var stdOut bytes.Buffer
	var stdErr bytes.Buffer
	cmd.Stdout = &stdOut
	cmd.Stderr = &stdErr

	if err := cmd.Run(); err != nil {
		// 读取标准错误输出
		return "", fmt.Errorf("执行 virsh 命令失败: %v, %s", err, stdErr.String())
	}
	return stdOut.String(), nil
}

type ListType string

const (
	Active   ListType = "active"
	Inactive ListType = "inactive"
	All      ListType = "all"
)

// GetVMList 获取虚拟机列表
func GetVMList(listType ListType) ([]VMInfo, error) {
	args := []string{"list"}

	switch listType {
	case Inactive:
		args = append(args, "--inactive")
	case All:
		args = append(args, "--all")
		// Active 是默认值，不需要额外参数
	}

	output, err := ExecVirshCommand(args...)
	if err != nil {
		return nil, err
	}

	// 解析文本输出
	return parseTextVMList(output)
}

type Disk struct {
	Type   string `xml:"type,attr"`
	Device string `xml:"device,attr"`
	Source struct {
		File string `xml:"file,attr"`
		Size string `xml:"size,attr"`
	} `xml:"source"`
	Driver struct {
		Name    string `xml:"name,attr"`
		Type    string `xml:"type,attr"`
		Discard string `xml:"discard,attr"`
	} `xml:"driver"`
}

type Domain struct {
	Disks    []Disk `xml:"devices>disk"`
	Graphics []struct {
		Type   string `xml:"type,attr"`
		Port   string `xml:"port,attr"`
		Listen string `xml:"listen,attr"`
	} `xml:"devices>graphics"`
	Hostdevs []struct {
		Type   string `xml:"type,attr"`
		Source struct {
			Address struct {
				Domain   string `xml:"domain,attr"`
				Bus      string `xml:"bus,attr"`
				Slot     string `xml:"slot,attr"`
				Function string `xml:"function,attr"`
			} `xml:"address"`
		} `xml:"source"`
	} `xml:"devices>hostdev"`
}

// extractDisk 解析xml内容 读取xpath=//disk@type='file' 且//disk@device='disk' 下的source/@file属性值
func extractDisk(xmlConfig string) []Disk {
	var domain Domain
	err := xml.Unmarshal([]byte(xmlConfig), &domain)
	if err != nil {
		return nil
	}
	var diskPaths []Disk
	for _, disk := range domain.Disks {
		if disk.Type == "file" && disk.Device == "disk" && disk.Source.File != "" {
			diskPaths = append(diskPaths, disk)
		}
	}
	return diskPaths
}

// CreateDiskFile 使用qemu-img命令创建磁盘文件
func CreateDiskFile(diskPath, format, size string) (string, error) {
	log.Printf("正在分配磁盘文件: %s, 格式: %s", diskPath, format)
	cmd := exec.Command("qemu-img", "create", "-f", format, diskPath, size)
	var stdOut bytes.Buffer
	var stdErr bytes.Buffer
	cmd.Stdout = &stdOut
	cmd.Stderr = &stdErr
	if err := cmd.Run(); err != nil {
		// 读取标准错误输出
		return "", fmt.Errorf("执行 qemu-img 命令失败: %v, %s", err, stdErr.String())
	}
	return stdOut.String(), nil
}

// parseTextVMList 解析 virsh list 命令的文本输出
func parseTextVMList(output string) ([]VMInfo, error) {
	lines := strings.Split(output, "\n")
	var vmList []VMInfo

	// 跳过标题行（前2行）
	for i := 2; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		if line == "" {
			continue
		}

		// 分割行内容，最多分割为3部分
		parts := strings.Fields(line)
		if len(parts) < 3 {
			continue
		}

		// 处理 ID, Name 和 State
		id := parts[0]
		name := parts[1]
		state := strings.Join(parts[2:], " ") // 状态可能包含空格

		vmList = append(vmList, VMInfo{
			ID:    id,
			Name:  name,
			State: state,
		})
	}

	return vmList, nil
}

// CheckAndActivateDefaultNetwork 检查并激活默认网络
func CheckAndActivateDefaultNetwork() error {
	// 检查网络状态（默认只显示活动网络）
	status, err := ExecVirshCommand("net-list", "--name")
	if err != nil {
		return fmt.Errorf("检查网络状态失败: %v", err)
	}

	// 检查默认网络是否已经激活
	if strings.Contains(status, "default") {
		return nil // 默认网络已激活
	}

	// 检查默认网络是否存在（包括未激活的）
	allNetworks, err := ExecVirshCommand("net-list", "--all", "--name")
	if err != nil {
		return fmt.Errorf("检查所有网络失败: %v", err)
	}

	if !strings.Contains(allNetworks, "default") {
		// 默认网络不存在，使用配置文件定义它
		_, err = ExecVirshCommand("net-define", "/etc/libvirt/qemu/networks/default.xml")
		if err != nil {
			return fmt.Errorf("定义默认网络失败: %v", err)
		}
	}

	// 尝试激活默认网络
	_, err = ExecVirshCommand("net-start", "default")
	if err != nil {
		return fmt.Errorf("激活默认网络失败: %v", err)
	}

	return nil
}

// StartVM 启动指定的虚拟机
func StartVM(vmName string) error {
	_, err := ExecVirshCommand("start", vmName)
	return err
}

// ShutdownVM 关闭指定的虚拟机
func ShutdownVM(vmName string) error {
	_, err := ExecVirshCommand("shutdown", vmName)
	return err
}

// SuspendVM 挂起指定的虚拟机
func SuspendVM(vmName string) error {
	_, err := ExecVirshCommand("suspend", vmName)
	return err
}

// ResumeVM 恢复挂起的虚拟机
func ResumeVM(vmName string) error {
	_, err := ExecVirshCommand("resume", vmName)
	return err
}

// GetVMInfo 获取指定虚拟机的详细信息
func GetVMInfo(vmName string) (string, error) {
	command, err := ExecVirshCommand("dumpxml", vmName)
	return command, err
}

// ForceShutdownVM 强制关闭指定的虚拟机
func ForceShutdownVM(vmName string) error {
	_, err := ExecVirshCommand("destroy", vmName)
	return err
}

// DeleteVM 删除指定的虚拟机
func DeleteVM(vmName string) error {
	_, err := ExecVirshCommand("undefine", vmName, "--nvram")
	return err
}

// SystemResourceInfo 表示系统资源信息
type SystemResourceInfo struct {
	CPUCores    int   `json:"cpuCores"`
	TotalMemory int64 `json:"totalMemory"` // MB
}

// GetSystemResourceInfo 获取系统CPU核心数量和内存大小信息
func GetSystemResourceInfo() (*SystemResourceInfo, error) {
	// 获取CPU核心数量
	cpuOutput, err := exec.Command("lscpu", "--parse=CPU").Output()
	if err != nil {
		return nil, fmt.Errorf("获取CPU信息失败: %v", err)
	}

	cpuLines := strings.Split(string(cpuOutput), "\n")
	cpuCores := 0
	for _, line := range cpuLines {
		line = strings.TrimSpace(line)
		if line != "" && !strings.HasPrefix(line, "#") {
			cpuCores++
		}
	}

	// 获取内存大小
	memOutput, err := exec.Command("free", "-m").Output()
	if err != nil {
		return nil, fmt.Errorf("获取内存信息失败: %v", err)
	}

	memLines := strings.Split(string(memOutput), "\n")
	var totalMemory int64
	if len(memLines) >= 2 {
		memParts := strings.Fields(memLines[1])
		if len(memParts) >= 2 {
			fmt.Sscanf(memParts[1], "%d", &totalMemory)
		}
	}

	return &SystemResourceInfo{
		CPUCores:    cpuCores,
		TotalMemory: totalMemory,
	}, nil
}

// CreateVMFromXML 通过xml创建虚拟机
func CreateVMFromXML(vmName, xmlConfig string) error {
	// 创建虚拟机前检查并激活默认网络
	if err := CheckAndActivateDefaultNetwork(); err != nil {
		return err
	}

	// 将xml文本写入到临时文件中,
	xmlPath := "/tmp/" + vmName + ".xml"
	err := os.WriteFile(xmlPath, []byte(xmlConfig), 0655)
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			log.Printf("Failed to remove temp file %s: %v", name, err)
		}
	}(xmlPath) // 创建完成后删除临时文件
	// 判断xml内定义的磁盘是否存在
	diskPaths := extractDisk(xmlConfig)
	for _, diskPath := range diskPaths {
		if _, err := os.Stat(diskPath.Source.File); os.IsNotExist(err) {
			// 创建磁盘文件
			_, err := CreateDiskFile(diskPath.Source.File, diskPath.Driver.Type, diskPath.Source.Size)
			if err != nil {
				return err
			}
		}
	}
	_, err = ExecVirshCommand("create", "--file", xmlPath)
	return err
}

// UpdateVMFromXML 通过xml更新虚拟机配置
func UpdateVMFromXML(vmName, xmlConfig string) error {
	// 将xml文本写入到临时文件中
	xmlPath := "/tmp/" + vmName + "_update.xml"
	err := os.WriteFile(xmlPath, []byte(xmlConfig), 0655)
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			log.Printf("Failed to remove temp file %s: %v", name, err)
		}
	}(xmlPath) // 更新完成后删除临时文件

	// 使用virsh define命令更新虚拟机配置
	_, err = ExecVirshCommand("define", "--file", xmlPath)
	return err
}

// AttachUsbDevice 为虚拟机添加usb设备
func AttachUsbDevice(vmName, usbId string) error {
	return editVmUsb(vmName, usbId, "attach-device")

}

// DetachUsbDevice 为虚拟机移除usb设备
// vmName 虚拟机名称
// usbId usb设备id如:0930:6545
func DetachUsbDevice(vmName, usbId string) error {
	return editVmUsb(vmName, usbId, "detach-device")
}

// AttachPCIDevice 为虚拟机添加PCI设备
func AttachPCIDevice(vmName, pciID string) error {
	return editVmPCI(vmName, pciID, "attach-device")
}

// DetachPCIDevice 为虚拟机移除PCI设备
func DetachPCIDevice(vmName, pciID string) error {
	return editVmPCI(vmName, pciID, "detach-device")
}

func editVmUsb(vmName, usbId, action string) error {
	xmlPath := fmt.Sprintf("/tmp/%s_%s_usb.xml", vmName, usbId)
	split := strings.Split(usbId, ":")
	vendor := split[0]
	device := split[1]
	xml := `
	<hostdev mode="subsystem" type="usb" managed="yes">
	  <source>
		<vendor id="0x%s"/>
		<product id="0x%s"/>
	  </source>
	  <alias name="hostdev0"/>
	  <address type="usb" bus="0" port="1"/>
	</hostdev>
	`
	xml = fmt.Sprintf(xml, vendor, device)
	err := os.WriteFile(xmlPath, []byte(xml), 0655)
	defer func() {
		err := os.Remove(xmlPath)
		if err != nil {
			log.Printf("Failed to remove temp file %s: %v", xmlPath, err)
		}
	}()
	if err != nil {
		return err
	}
	_, err = ExecVirshCommand(action, vmName, xmlPath)
	if err != nil {
		return err
	}
	return nil
}

func editVmPCI(vmName, pciID, action string) error {
	xmlPath := fmt.Sprintf("/tmp/%s_%s_pci.xml", vmName, pciID)
	xml := GeneratePCIHostdevXML(pciID)
	err := os.WriteFile(xmlPath, []byte(xml), 0655)
	defer func() {
		err := os.Remove(xmlPath)
		if err != nil {
			log.Printf("Failed to remove temp file %s: %v", xmlPath, err)
		}
	}()
	if err != nil {
		return err
	}
	_, err = ExecVirshCommand(action, vmName, xmlPath)
	if err != nil {
		return err
	}
	return nil
}
