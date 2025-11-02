package vmopts

import "testing"

// 测试获取所有虚拟机列表
func TestGetVMList_All(t *testing.T) {
	vms, err := GetVMList(All)
	if err != nil {
		t.Errorf("GetVMList(All) failed: %v", err)
	}
	if len(vms) == 0 {
		t.Errorf("GetVMList(All) returned empty list")
	}
}

// TestGetVMInfo 测试获取虚拟机信息
func TestGetVMInfo(t *testing.T) {
	vmName := "ubuntu2510"
	info, err := GetVMInfo(vmName)
	if err != nil {
		t.Errorf("GetVMInfo(%s) failed: %v", vmName, err)
	}
	if info == "" {
		t.Errorf("GetVMInfo(%s) returned empty info", vmName)
	}
}

// StartVM 测试
func TestStartVM(t *testing.T) {
	vmName := "ubuntu2510"
	err := StartVM(vmName)
	if err != nil {
		t.Errorf("StartVM(%s) failed: %v", vmName, err)
	}
}

// ShutdownVM 测试
func TestShutdownVM(t *testing.T) {
	vmName := "ubuntu2510"
	err := ShutdownVM(vmName)
	if err != nil {
		t.Errorf("ShutdownVM(%s) failed: %v", vmName, err)
	}
}

// SuspendVM 测试
func TestSuspendVM(t *testing.T) {
	vmName := "ubuntu2510"
	err := SuspendVM(vmName)
	if err != nil {
		t.Errorf("SuspendVM(%s) failed: %v", vmName, err)
	}
}

// ResumeVM 测试
func TestResumeVM(t *testing.T) {
	vmName := "ubuntu2510"
	err := ResumeVM(vmName)
	if err != nil {
		t.Errorf("ResumeVM(%s) failed: %v", vmName, err)
	}
}

// ForceShutdownVM 测试
func TestForceShutdownVM(t *testing.T) {
	vmName := "ubuntu2510"
	err := ForceShutdownVM(vmName)
	if err != nil {
		t.Errorf("ForceShutdownVM(%s) failed: %v", vmName, err)
	}
}
