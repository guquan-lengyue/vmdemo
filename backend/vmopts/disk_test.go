package vmopts

import (
	"testing"
)

func TestDiskPool(t *testing.T) {
	poolPath := "./test_disk_pool"
	diskName := "test_disk.qcow2"
	diskFormat := "qcow2"
	diskSize := "10G"
	newDiskSize := "20G"

	// 初始化硬盘池
	t.Log("Initializing disk pool...")
	pool, err := NewDiskPool(poolPath)
	if err != nil {
		t.Fatalf("Failed to initialize disk pool: %s", err)
	}

	// 添加虚拟硬盘
	t.Log("Adding disk to pool...")
	err = pool.AddDisk(diskName, diskFormat, diskSize)
	if err != nil {
		t.Fatalf("Failed to add disk to pool: %s", err)
	}

	// 列出硬盘
	t.Log("Listing disks in pool...")
	disks, err := pool.ListDisks()
	if err != nil {
		t.Fatalf("Failed to list disks in pool: %s", err)
	}
	t.Logf("Disks in pool: %v", disks)

	// 获取硬盘信息
	t.Log("Getting disk info...")
	info, err := pool.GetDiskInfo(diskName)
	if err != nil {
		t.Fatalf("Failed to get disk info: %s", err)
	}
	t.Logf("Disk Info: %+v", info)

	// 调整硬盘大小
	t.Log("Resizing disk...")
	err = pool.ResizeDisk(diskName, newDiskSize)
	if err != nil {
		t.Fatalf("Failed to resize disk: %s", err)
	}

	// 删除硬盘
	t.Log("Deleting disk from pool...")
	err = pool.DeleteDisk(diskName)
	if err != nil {
		t.Fatalf("Failed to delete disk from pool: %s", err)
	}
}
