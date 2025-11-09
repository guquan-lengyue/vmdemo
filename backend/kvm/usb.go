package kvm

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"
)

type UsbInfo struct {
	ID   string `json:"id"`
	Info string `json:"info"`
}

var idReg = regexp.MustCompile(`ID\s([0-9a-f]{4}:[0-9a-f]{4})`)

// GetUsbList 读取usb列表
func GetUsbList() ([]UsbInfo, error) {
	cmd := exec.Command("lsusb")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to get usb list: %s, error: %s", string(output), err)
	}
	lines := strings.Split(string(output), "\n")
	var usbList []UsbInfo
	for _, line := range lines {
		if line == "" {
			continue
		}
		id := idReg.FindStringSubmatch(line)
		if len(id) == 0 {
			continue
		}
		usbList = append(usbList, UsbInfo{
			ID:   id[1],
			Info: line,
		})
	}
	return usbList, nil
}
