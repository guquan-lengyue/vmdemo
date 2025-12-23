package main

import (
	"encoding/json"
	"fmt"
)

// TestStruct 简单的测试结构体
type TestStruct struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func main() {
	// 创建测试数据
	testData := []TestStruct{
		{ID: "1", Name: "Test1"},
		{ID: "2", Name: "Test2"},
	}

	// 序列化为JSON
	jsonData, err := json.Marshal(testData)
	if err != nil {
		fmt.Printf("Error marshaling to JSON: %v\n", err)
		return
	}

	// 打印JSON
	fmt.Printf("JSON output: %s\n\n", string(jsonData))

	// 格式化JSON
	formattedJSON, err := json.MarshalIndent(testData, "", "  ")
	if err != nil {
		fmt.Printf("Error formatting JSON: %v\n", err)
		return
	}

	fmt.Printf("Formatted JSON: %s\n", string(formattedJSON))
}
