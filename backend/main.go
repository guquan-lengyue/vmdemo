package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"vmdemo/service"
)

func main() {
	// 启动gin接口服务
	r := gin.Default()
	service.RegisterVMOpsRoutes(r)
	err := r.Run(":8080")
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
