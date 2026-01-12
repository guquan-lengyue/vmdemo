package service

import (
	"net/http"
	"vmdemo/kvm"

	"github.com/gin-gonic/gin"
)

// GetPCIList 注册PCI设备相关的路由
func GetPCIList(r *gin.RouterGroup) {
	g := r.Group("/pci")
	{
		// 获取PCI设备列表
		g.GET("/list", func(c *gin.Context) {
			pciList, err := kvm.GetPCIList()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
				return
			}
			c.JSON(http.StatusOK, pciList)
		})

		// 获取指定PCI设备的详细信息
		g.GET("/detail", func(c *gin.Context) {
			pciID := c.Query("id")
			if pciID == "" {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "PCI device ID is required",
				})
				return
			}

			pciDetail, err := kvm.GetPCIDeviceDetail(pciID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
				return
			}

			c.JSON(http.StatusOK, pciDetail)
		})
	}
}
