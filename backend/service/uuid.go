package service

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// GenerateUUIDHandler 生成UUID的接口
func GenerateUUIDHandler(c *gin.Context) {
	// 生成V4版本的UUID
	generatedUUID := uuid.New().String()

	// 直接构造JSON字符串
	response := fmt.Sprintf(`{"uuid":"%s","version":"v4"}`, generatedUUID)
	c.Data(http.StatusOK, "application/json", []byte(response))
}

// RegisterUUIDRoutes 注册UUID生成的路由
func RegisterUUIDRoutes(router *gin.RouterGroup) {
	uuidGroup := router.Group("/uuid")
	{
		uuidGroup.GET("/generate", GenerateUUIDHandler)
	}
}
