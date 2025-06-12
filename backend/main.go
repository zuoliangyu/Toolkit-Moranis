package main

import (
	"log"
	"servermanager/internal/config"
	"servermanager/internal/handlers"
	"servermanager/internal/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	// 先加载配置
	if err := config.InitConfig(); err != nil {
		log.Fatalf("配置加载失败: %v", err)
	}

	// 初始化数据库连接
	db, err := config.InitDB()
	if err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}

	// 创建处理器
	publicHandler := handlers.NewPublicHandler(db)
	adminHandler := handlers.NewAdminHandler(db)

	// 创建公共服务器
	publicServer := gin.New()
	publicServer.Use(gin.Logger())
	publicServer.Use(gin.Recovery())
	publicServer.Use(middleware.CORSMiddleware())
	// 静态资源服务
	publicServer.Static("/uploads", "./uploads")

	// 公共API路由组
	publicAPI := publicServer.Group("/api")
	{
		publicAPI.GET("/materials", publicHandler.GetMaterials)
		publicAPI.GET("/materials/:id/download", publicHandler.DownloadMaterial)
		publicAPI.GET("/folders", publicHandler.GetFolders)
		publicAPI.GET("/folders/:id/materials", publicHandler.GetFolderMaterials)
	}

	// 创建管理员服务器
	adminServer := gin.New()
	adminServer.Use(gin.Logger())
	adminServer.Use(gin.Recovery())
	adminServer.Use(middleware.CORSMiddleware())

	publicServer.Static("/uploads", "./uploads")

	// 管理员API路由组
	adminAPI := adminServer.Group("/api")
	{
		// 不需要认证的路由
		adminAPI.POST("/login", adminHandler.Login)
		adminAPI.POST("/admin/login", adminHandler.Login)

		// 需要认证的路由
		adminAuth := adminAPI.Group("")
		adminAuth.Use(middleware.AuthMiddleware())
		{
			adminAuth.POST("/materials", adminHandler.UploadMaterial)
			adminAuth.DELETE("/materials/:id", adminHandler.DeleteMaterial)
			adminAuth.GET("/folders", adminHandler.GetFolders)
			adminAuth.POST("/folders", adminHandler.CreateFolder)
			adminAuth.PUT("/folders/:id", adminHandler.UpdateFolder)
			adminAuth.DELETE("/folders/:id", adminHandler.DeleteFolder)
		}
	}

	// 启动服务器
	go func() {
		log.Printf("公共服务器启动在端口 8080")
		if err := publicServer.Run(":8080"); err != nil {
			log.Fatalf("公共服务器启动失败: %v", err)
		}
	}()

	log.Printf("管理员服务器启动在端口 8081")
	if err := adminServer.Run(":8081"); err != nil {
		log.Fatalf("管理员服务器启动失败: %v", err)
	}
}
