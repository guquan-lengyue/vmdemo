// Package vmopts virsh 命令调用工具
package vmopts

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
	Disks []Disk `xml:"devices>disk"`
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
	return ExecVirshCommand("dominfo", vmName)
}

// ForceShutdownVM 强制关闭指定的虚拟机
func ForceShutdownVM(vmName string) error {
	_, err := ExecVirshCommand("destroy", vmName)
	return err
}

// DeleteVM 删除指定的虚拟机
func DeleteVM(vmName string) error {
	_, err := ExecVirshCommand("undefine", vmName, "--nvram", "--remove-all-storage")
	return err
}

// CreateVMFromXML 通过xml创建虚拟机
func CreateVMFromXML(vmName, xmlConfig string) error {
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
