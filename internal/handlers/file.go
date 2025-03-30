package handlers

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func UploadFile(c *gin.Context) {
	// Get user from context (set by auth middleware)
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	// Get file from form data
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create uploads directory if not exists
	if err := os.MkdirAll("uploads", os.ModePerm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create upload directory"})
		return
	}

	// Generate unique filename
	fileID := uuid.New().String()
	fileExt := filepath.Ext(file.Filename)
	newFilename := fileID + fileExt
	filePath := filepath.Join("uploads", newFilename)

	// Save file locally
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save file"})
		return
	}

	// Save file metadata to database
	dbFile, err := database.CreateFile(database.File{
		ID:       fileID,
		UserID:   user.(*database.User).ID,
		Name:     file.Filename,
		Path:     filePath,
		Size:     file.Size,
		MimeType: file.Header.Get("Content-Type"),
	})
	if err != nil {
		// Clean up saved file if DB operation fails
		os.Remove(filePath)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save file metadata"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "file uploaded successfully",
		"file":    dbFile,
		"url":     "/files/" + fileID,
	})
}

func GetUserFiles(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	files, err := database.GetFilesByUserID(user.(*database.User).ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get files"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"files": files,
	})
}

func DownloadFile(c *gin.Context) {
	fileID := c.Param("id")

	file, err := database.GetFileByID(fileID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "file not found"})
		return
	}

	// Check if file is public or belongs to requesting user
	user, _ := c.Get("user")
	if !file.IsPublic && file.UserID != user.(*database.User).ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "unauthorized access"})
		return
	}

	c.File(file.Path)
}
