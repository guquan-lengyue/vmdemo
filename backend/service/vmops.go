package service

import (
	"net/http"
	"vmdemo/vmopts"

	"github.com/gin-gonic/gin"
)

// StartVMHandler 启动虚拟机的接口
func StartVMHandler(c *gin.Context) {
	vmName := c.Query("name")
	if vmName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "VM name is required"})
		return
	}

	err := vmopts.StartVM(vmName)
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

	err := vmopts.ShutdownVM(vmName)
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

	err := vmopts.SuspendVM(vmName)
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

	err := vmopts.ResumeVM(vmName)
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

	err := vmopts.ForceShutdownVM(vmName)
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

	err := vmopts.DeleteVM(vmName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "VM deleted successfully", "vmName": vmName})
}

// ListVMsHandler 获取虚拟机列表的接口
func ListVMsHandler(c *gin.Context) {
	listType := c.DefaultQuery("type", "all") // 默认获取所有虚拟机
	var vmListType vmopts.ListType

	switch listType {
	case "active":
		vmListType = vmopts.Active
	case "inactive":
		vmListType = vmopts.Inactive
	case "all":
		vmListType = vmopts.All
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid list type"})
		return
	}

	vms, err := vmopts.GetVMList(vmListType)
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

	info, err := vmopts.GetVMInfo(vmName)
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

	err := vmopts.CreateVMFromXML(vmName, xmlConfig)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "VM created successfully", "vmName": vmName})
}

// RegisterVMOpsRoutes 注册虚拟机操作的路由
func RegisterVMOpsRoutes(router *gin.Engine) {
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
	}
}
