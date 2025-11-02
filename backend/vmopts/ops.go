// Package vmopts virsh 命令调用工具
package vmopts

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// VMInfo 表示虚拟机的基本信息
type VMInfo struct {
	ID    string
	Name  string
	State string
}

// ExecVirshCommand 执行 virsh 命令并返回输出结果
func ExecVirshCommand(args ...string) (string, error) {
	cmd := exec.Command("virsh", args...)
	var out bytes.Buffer
	var err bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &err

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("执行 virsh 命令失败: %v, %s", err, err.Error())
	}

	return out.String(), nil
}

// GetVMList 读取虚拟机列表
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

// parseTextVMList 解析 virsh list 命令的文本输出
func parseTextVMList(output string) ([]VMInfo, error) {
	lines := strings.Split(output, "\n")
	var vmList []VMInfo

	// 跳过标题行（前2行）
	for i := 2; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		if line == "" || strings.HasPrefix(line, "-") {
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
