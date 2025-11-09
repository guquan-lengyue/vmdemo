package main

import (
	"log"
	"vmdemo/service"

	"github.com/gin-gonic/gin"
)

const (
	DefaultDiskPoolPath = "./disks_pool"
)

func main() {
	// 启动gin接口服务
	r := gin.Default()
	g := r.Group("/api")
	service.RegisterVMOpsRoutes(g)
	err := service.RegisterDiskRoutes(g, DefaultDiskPoolPath)
	if err != nil {
		log.Fatalf("Error registering disk routes: %v", err)
	}
	service.CrateVncWs(g)
	err = r.Run(":8080")
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
