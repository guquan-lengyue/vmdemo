package main

import (
	"encoding/json"
	"fmt"
	"vmdemo/kvm"
)

func main() {
	// 获取PCI设备列表
	pciList, err := kvm.GetPCIList()
	if err != nil {
		fmt.Printf("Error getting PCI list: %v\n", err)
		return
	}

	// 打印原始列表
	fmt.Printf("Raw PCI list: %+v\n\n", pciList)

	// 序列化为JSON
	jsonData, err := json.Marshal(pciList)
	if err != nil {
		fmt.Printf("Error marshaling to JSON: %v\n", err)
		return
	}

	// 打印JSON
	fmt.Printf("JSON output: %s\n\n", string(jsonData))

	// 格式化JSON
	formattedJSON, err := json.MarshalIndent(pciList, "", "  ")
	if err != nil {
		fmt.Printf("Error formatting JSON: %v\n", err)
		return
	}

	fmt.Printf("Formatted JSON: %s\n", string(formattedJSON))
}
