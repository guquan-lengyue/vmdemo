package service

import (
	"net/http"
	"vmdemo/kvm"

	"github.com/gin-gonic/gin"
)

// StartVMHandler 启动虚拟机的接口
func StartVMHandler(c *gin.Context) {
	vmName := c.Query("name")
	if vmName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "VM name is required"})
		return
	}

	err := kvm.StartVM(vmName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "VM started successfully", "vmName": vmName})
}

// StopVMHandler 停止虚拟机的接口
func StopVMHandler(c *gin.Context) {
	vmName := c.Query("name")
	if vmName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "VM name is required"})
		return
	}

	err := kvm.ShutdownVM(vmName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "VM stopped successfully", "vmName": vmName})
}

// SuspendVMHandler 挂起虚拟机的接口
func SuspendVMHandler(c *gin.Context) {
	vmName := c.Query("name")
	if vmName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "VM name is required"})
		return
	}

	err := kvm.SuspendVM(vmName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "VM suspended successfully", "vmName": vmName})
}

// ResumeVMHandler 恢复虚拟机的接口
func ResumeVMHandler(c *gin.Context) {
	vmName := c.Query("name")
	if vmName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "VM name is required"})
		return
	}

	err := kvm.ResumeVM(vmName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "VM resumed successfully", "vmName": vmName})
}

// ForceShutdownVMHandler 强制关闭虚拟机的接口
func ForceShutdownVMHandler(c *gin.Context) {
	vmName := c.Query("name")
	if vmName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "VM name is required"})
		return
	}

	err := kvm.ForceShutdownVM(vmName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "VM forcefully shut down", "vmName": vmName})
}

// DeleteVMHandler 删除虚拟机的接口
func DeleteVMHandler(c *gin.Context) {
	vmName := c.Query("name")
	if vmName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "VM name is required"})
		return
	}

	err := kvm.DeleteVM(vmName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "VM deleted successfully", "vmName": vmName})
}

// ListVMsHandler 获取虚拟机列表的接口
func ListVMsHandler(c *gin.Context) {
	listType := c.DefaultQuery("type", "all") // 默认获取所有虚拟机
	var vmListType kvm.ListType

	switch listType {
	case "active":
		vmListType = kvm.Active
	case "inactive":
		vmListType = kvm.Inactive
	case "all":
		vmListType = kvm.All
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid list type"})
		return
	}

	vms, err := kvm.GetVMList(vmListType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"vms": vms})
}

// GetVMInfoHandler 获取虚拟机详细信息的接口
func GetVMInfoHandler(c *gin.Context) {
	vmName := c.Query("name")
	if vmName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "VM name is required"})
		return
	}

	info, err := kvm.GetVMInfo(vmName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"vmInfo": info})
}

// CreateVMHandler 通过 XML 创建虚拟机的接口
func CreateVMHandler(c *gin.Context) {
	vmName := c.PostForm("name")
	xmlConfig := c.PostForm("xmlConfig")

	if vmName == "" || xmlConfig == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "VM name and XML config are required"})
		return
	}

	err := kvm.CreateVMFromXML(vmName, xmlConfig)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "VM created successfully", "vmName": vmName})
}

// AttachUsbDeviceHandler 为虚拟机添加usb设备的接口
func AttachUsbDeviceHandler(c *gin.Context) {
	vmName := c.Query("name")
	usbId := c.Query("usbId")

	if vmName == "" || usbId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "VM name and USB XML are required"})
		return
	}

	err := kvm.AttachUsbDevice(vmName, usbId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "USB device attached successfully", "vmName": vmName})
}

// DetachUsbDeviceHandler 为虚拟机移除usb设备的接口
func DetachUsbDeviceHandler(c *gin.Context) {
	vmName := c.Query("name")
	usbId := c.Query("usbId")

	if vmName == "" || usbId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "VM name and USB ID are required"})
		return
	}

	err := kvm.DetachUsbDevice(vmName, usbId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "USB device detached successfully", "vmName": vmName})
}

// RegisterVMOpsRoutes 注册虚拟机操作的路由
func RegisterVMOpsRoutes(router *gin.RouterGroup) {
	vmGroup := router.Group("/vm")
	{
		vmGroup.GET("/start", StartVMHandler)
		vmGroup.GET("/stop", StopVMHandler)
		vmGroup.GET("/suspend", SuspendVMHandler)
		vmGroup.GET("/resume", ResumeVMHandler)
		vmGroup.GET("/force-shutdown", ForceShutdownVMHandler)
		vmGroup.GET("/delete", DeleteVMHandler)
		vmGroup.GET("/list", ListVMsHandler)
		vmGroup.GET("/info", GetVMInfoHandler)
		vmGroup.POST("/create", CreateVMHandler)
		vmGroup.GET("/attach-usb", AttachUsbDeviceHandler)
		vmGroup.GET("/detach-usb", DetachUsbDeviceHandler)
	}
}
