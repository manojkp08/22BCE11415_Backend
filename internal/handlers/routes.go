package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/manojkp08/22BCE11415_Backend/pkg/middleware"
)

func SetupRoutes(router *gin.Engine) {
	// Auth routes
	router.GET("/auth/google/login", GoogleLoginHandler)
	router.GET("/auth/google/callback", GoogleCallbackHandler)

	// File routes (protected with JWT auth)
	authGroup := router.Group("/")
	authGroup.Use(middleware.AuthMiddleware())
	{
		authGroup.POST("/upload", UploadFile)
		authGroup.GET("/files", GetUserFiles)
		authGroup.GET("/files/:id", DownloadFile)
	}
}
