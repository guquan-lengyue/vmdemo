package kvm

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"
)

// PCIInfo 表示PCI设备的基本信息
type PCIInfo struct {
	ID     string `json:"id"`     // PCI设备ID，格式如：0000:00:02.0
	Name   string `json:"name"`   // PCI设备名称
	Class  string `json:"class"`  // PCI设备类别
	Vendor string `json:"vendor"` // PCI设备厂商
	Device string `json:"device"` // PCI设备型号
}

// GetPCIList 获取系统中所有PCI设备列表
func GetPCIList() ([]PCIInfo, error) {
	cmd := exec.Command("lspci", "-nn")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to get PCI list: %s, error: %s", string(output), err)
	}
	return parsePCIOutput(string(output)), nil
}

// GetPCIDeviceDetail 获取指定PCI设备的详细信息
func GetPCIDeviceDetail(pciID string) (*PCIInfo, error) {
	// 去掉域名前缀（如果有），因为lspci命令不识别带域名的格式
	if strings.HasPrefix(pciID, "0000:") {
		pciID = pciID[5:]
	}

	// 使用-s选项指定设备地址
	cmd := exec.Command("lspci", "-nn", "-s", pciID)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to get PCI device detail: %s, error: %s", string(output), err)
	}
	pciList := parsePCIOutput(string(output))
	if len(pciList) > 0 {
		return &pciList[0], nil
	}
	return nil, fmt.Errorf("PCI device not found: %s", pciID)
}

// parsePCIOutput 解析lspci命令的输出
func parsePCIOutput(output string) []PCIInfo {
	lines := strings.Split(output, "\n")
	var pciList []PCIInfo

	// 正则表达式用于解析lspci输出行
	// 示例行：00:02.0 VGA compatible controller [0300]: Intel Corporation Device [8086:4626] (rev 0c)
	lineRegex := regexp.MustCompile(`^([0-9a-fA-F:]+\.[0-9a-fA-F])\s+(.+)\s+\[([0-9a-fA-F]{4})\]:\s+(.+)\s+\[([0-9a-fA-F]{4}):([0-9a-fA-F]{4})\]`)

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		matches := lineRegex.FindStringSubmatch(line)
		if len(matches) < 7 {
			continue
		}

		pciInfo := PCIInfo{
			ID:     fmt.Sprintf("0000:%s", matches[1]), // 添加域名前缀
			Name:   matches[2],
			Class:  matches[3],
			Vendor: matches[5],
			Device: matches[6],
		}

		pciList = append(pciList, pciInfo)
	}

	return pciList
}

// GeneratePCIHostdevXML 生成PCI设备的hostdev XML配置
func GeneratePCIHostdevXML(pciID string) string {
	return fmt.Sprintf(`
	<hostdev mode="subsystem" type="pci" managed="yes">
	  <source>
	    <address domain="0x0000" bus="0x%s" slot="0x%s" function="0x%s"/>
	  </source>
	  <alias name="hostdev0"/>
	  <address type="pci" domain="0x0000" bus="0x00" slot="0x05" function="0x0" multifunction="on"/>
	</hostdev>`, getPCIBus(pciID), getPCISlot(pciID), getPCIFunction(pciID))
}

// 辅助函数：从PCI ID中提取总线号
func getPCIBus(pciID string) string {
	parts := strings.Split(pciID, ":")
	if len(parts) >= 2 {
		busSlot := parts[len(parts)-1]
		bus := strings.Split(busSlot, ".")[0]
		return bus
	}
	return "00"
}

// 辅助函数：从PCI ID中提取槽位号
func getPCISlot(pciID string) string {
	parts := strings.Split(pciID, ":")
	if len(parts) >= 2 {
		busSlot := parts[len(parts)-1]
		busSlotParts := strings.Split(busSlot, ".")[0]
		slot := busSlotParts[2:]
		return slot
	}
	return "00"
}

// 辅助函数：从PCI ID中提取功能号
func getPCIFunction(pciID string) string {
	parts := strings.Split(pciID, ".")
	if len(parts) >= 2 {
		return parts[len(parts)-1]
	}
	return "0"
}
