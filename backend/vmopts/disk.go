package vmopts

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// DiskInfo 表示虚拟盘的信息
type DiskInfo struct {
	Format string
	Size   string
	Path   string
}

// CreateDisk 创建虚拟盘
func CreateDisk(path, format, size string) error {
	cmd := exec.Command("qemu-img", "create", "-f", format, path, size)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to create disk: %s, error: %s", string(output), err)
	}
	return nil
}

// GetDiskInfo 获取虚拟盘信息
func GetDiskInfo(path string) (*DiskInfo, error) {
	cmd := exec.Command("qemu-img", "info", "--output=json", path)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to get disk info: %s, error: %s", string(output), err)
	}

	// 解析 JSON 输出（简化处理）
	info := &DiskInfo{}
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, "\"format\"") {
			info.Format = strings.TrimSpace(strings.Split(line, ":")[1])
		} else if strings.Contains(line, "\"virtual-size\"") {
			info.Size = strings.TrimSpace(strings.Split(line, ":")[1])
		} else if strings.Contains(line, "\"filename\"") {
			info.Path = strings.Trim(strings.Split(line, ":")[1], "\" ")
		}
	}
	return info, nil
}

// ResizeDisk 调整虚拟盘大小
func ResizeDisk(path, newSize string) error {
	cmd := exec.Command("qemu-img", "resize", path, newSize)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to resize disk: %s, error: %s", string(output), err)
	}
	return nil
}

// DeleteDisk 删除虚拟盘
func DeleteDisk(path string) error {
	cmd := exec.Command("rm", "-f", path)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to delete disk: %s, error: %s", string(output), err)
	}
	return nil
}

// DiskPool 表示虚拟硬盘池
type DiskPool struct {
	PoolPath string
}

// NewDiskPool 创建一个新的硬盘池
func NewDiskPool(poolPath string) (*DiskPool, error) {
	// 如果目录不存在，则创建
	if _, err := os.Stat(poolPath); os.IsNotExist(err) {
		err := os.MkdirAll(poolPath, 0755)
		if err != nil {
			return nil, fmt.Errorf("failed to create disk pool directory: %s", err)
		}
	}
	return &DiskPool{PoolPath: poolPath}, nil
}

// ListDisks 列出硬盘池中的所有虚拟硬盘
func (dp *DiskPool) ListDisks() ([]string, error) {
	files, err := os.ReadDir(dp.PoolPath)
	if err != nil {
		return nil, fmt.Errorf("failed to list disks in pool: %s", err)
	}

	var disks []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".qcow2") {
			disks = append(disks, file.Name())
		}
	}
	return disks, nil
}

// AddDisk 添加虚拟硬盘到硬盘池
func (dp *DiskPool) AddDisk(name, format, size string) error {
	diskPath := filepath.Join(dp.PoolPath, name)
	return CreateDisk(diskPath, format, size)
}

// DeleteDisk 从硬盘池中删除虚拟硬盘
func (dp *DiskPool) DeleteDisk(name string) error {
	diskPath := filepath.Join(dp.PoolPath, name)
	return DeleteDisk(diskPath)
}

// GetDiskInfo 获取硬盘池中某个虚拟硬盘的信息
func (dp *DiskPool) GetDiskInfo(name string) (*DiskInfo, error) {
	diskPath := filepath.Join(dp.PoolPath, name)
	return GetDiskInfo(diskPath)
}

// ResizeDisk 调整硬盘池中某个虚拟硬盘的大小
func (dp *DiskPool) ResizeDisk(name, newSize string) error {
	diskPath := filepath.Join(dp.PoolPath, name)
	return ResizeDisk(diskPath, newSize)
}
