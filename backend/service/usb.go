package service

import (
	"net/http"
	"vmdemo/kvm"

	"github.com/gin-gonic/gin"
)

func GetUsbList(r *gin.RouterGroup) {
	g := r.Group("/usb")
	{
		g.GET("/list", func(c *gin.Context) {
			usbList, err := kvm.GetUsbList()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
				return
			}
			c.JSON(http.StatusOK, usbList)
		})
	}
}
