package service

import (
	"net/http"
	"vmdemo/vmopts"

	"github.com/gin-gonic/gin"
)

// DiskPoolHandler 表示硬盘池的操作接口
type DiskPoolHandler struct {
	DiskPool *vmopts.DiskPool
}

// NewDiskPoolHandler 创建一个新的硬盘池处理器
func NewDiskPoolHandler(poolPath string) (*DiskPoolHandler, error) {
	diskPool, err := vmopts.NewDiskPool(poolPath)
	if err != nil {
		return nil, err
	}
	return &DiskPoolHandler{DiskPool: diskPool}, nil
}

// ListDisksHandler 列出硬盘池中的所有虚拟硬盘
func (h *DiskPoolHandler) ListDisksHandler(c *gin.Context) {
	disks, err := h.DiskPool.ListDisks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"disks": disks})
}

// AddDiskHandler 添加虚拟硬盘到硬盘池
func (h *DiskPoolHandler) AddDiskHandler(c *gin.Context) {
	var request struct {
		Name   string `json:"name" binding:"required"`
		Format string `json:"format" binding:"required"`
		Size   string `json:"size" binding:"required"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.DiskPool.AddDisk(request.Name, request.Format, request.Size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Disk created successfully"})
}

// GetDiskInfoHandler 获取硬盘池中某个虚拟硬盘的信息
func (h *DiskPoolHandler) GetDiskInfoHandler(c *gin.Context) {
	name := c.Param("name")
	info, err := h.DiskPool.GetDiskInfo(name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"diskInfo": info})
}

// ResizeDiskHandler 调整硬盘池中某个虚拟硬盘的大小
func (h *DiskPoolHandler) ResizeDiskHandler(c *gin.Context) {
	var request struct {
		Name    string `json:"name" binding:"required"`
		NewSize string `json:"newSize" binding:"required"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.DiskPool.ResizeDisk(request.Name, request.NewSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Disk resized successfully"})
}

// DeleteDiskHandler 从硬盘池中删除虚拟硬盘
func (h *DiskPoolHandler) DeleteDiskHandler(c *gin.Context) {
	name := c.Param("name")
	err := h.DiskPool.DeleteDisk(name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Disk deleted successfully"})
}

// RegisterDiskRoutes 注册硬盘池操作的路由
func RegisterDiskRoutes(router *gin.RouterGroup, poolPath string) error {
	handler, err := NewDiskPoolHandler(poolPath)
	if err != nil {
		return err
	}
	diskGroup := router.Group("/disk")
	{
		diskGroup.GET("/list", handler.ListDisksHandler)
		diskGroup.POST("/add", handler.AddDiskHandler)
		diskGroup.GET("/info/:name", handler.GetDiskInfoHandler)
		diskGroup.POST("/resize", handler.ResizeDiskHandler)
		diskGroup.DELETE("/delete/:name", handler.DeleteDiskHandler)
	}
	return nil
}
