package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/manojkp08/22BCE11415_Backend/internal/cache"
	"github.com/manojkp08/22BCE11415_Backend/internal/database"
)

// func UploadFile(c *gin.Context) {
// 	// Get user from context (set by auth middleware)
// 	user, exists := c.Get("user")
// 	if !exists {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
// 		return
// 	}

// 	// Get file from form data
// 	file, err := c.FormFile("file")
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	// Create uploads directory if not exists
// 	if err := os.MkdirAll("uploads", os.ModePerm); err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create upload directory"})
// 		return
// 	}

// 	// Generate unique filename
// 	fileID := uuid.New().String()
// 	fileExt := filepath.Ext(file.Filename)
// 	newFilename := fileID + fileExt
// 	filePath := filepath.Join("uploads", newFilename)

// 	// Save file locally
// 	if err := c.SaveUploadedFile(file, filePath); err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save file"})
// 		return
// 	}

// 	// Save file metadata to database
// 	dbFile, err := database.CreateFile(database.File{
// 		ID:       fileID,
// 		UserID:   user.(*database.User).ID,
// 		Name:     file.Filename,
// 		Path:     filePath,
// 		Size:     file.Size,
// 		MimeType: file.Header.Get("Content-Type"),
// 	})
// 	if err != nil {
// 		// Clean up saved file if DB operation fails
// 		os.Remove(filePath)
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save file metadata"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{
// 		"message": "file uploaded successfully",
// 		"file":    dbFile,
// 		"url":     "/files/" + fileID,
// 	})
// }

func UploadFile(c *gin.Context) {
	// Get user from context
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

	// Create channels for concurrent processing
	resultChan := make(chan *database.File)
	errorChan := make(chan error)
	done := make(chan bool)

	// Start concurrent file processing
	go func() {
		// Generate file metadata
		fileID := uuid.New().String()
		fileExt := filepath.Ext(file.Filename)
		newFilename := fileID + fileExt
		filePath := filepath.Join("uploads", newFilename)

		// Ensure upload directory exists
		if err := os.MkdirAll("uploads", os.ModePerm); err != nil {
			errorChan <- fmt.Errorf("failed to create upload directory: %w", err)
			return
		}

		// Save file to storage (local/S3)
		if err := saveUploadedFileConcurrently(file, filePath); err != nil {
			errorChan <- fmt.Errorf("failed to save file: %w", err)
			return
		}

		// Create file metadata
		dbFile := database.File{
			ID:       fileID,
			UserID:   user.(*database.User).ID,
			Name:     file.Filename,
			Path:     filePath,
			Size:     file.Size,
			MimeType: file.Header.Get("Content-Type"),
			IsPublic: false,
		}

		// Save to database
		createdFile, err := database.CreateFile(dbFile)
		if err != nil {
			// Clean up file if DB operation fails
			os.Remove(filePath)
			errorChan <- fmt.Errorf("failed to save metadata: %w", err)
			return
		}

		// Cache the file metadata
		fileJson, _ := json.Marshal(createdFile)
		if err := cache.SetFileMetadata(fileID, string(fileJson), 24*time.Hour); err != nil {
			log.Printf("Failed to cache file metadata: %v", err)
		}

		resultChan <- createdFile
	}()

	// Handle the concurrent operations
	go func() {
		select {
		case createdFile := <-resultChan:
			// Notify client via WebSocket
			websocket.BroadcastToUser(user.(*database.User).ID, gin.H{
				"event": "upload_complete",
				"file":  createdFile,
			})
			done <- true

		case err := <-errorChan:
			log.Printf("Upload failed: %v", err)
			errorChan <- err
		}
	}()

	// Respond immediately while processing continues in background
	select {
	case <-done:
		c.JSON(http.StatusOK, gin.H{
			"message": "File upload processed successfully",
			"status":  "completed",
		})
	case err := <-errorChan:
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  err.Error(),
			"status": "failed",
		})
	case <-time.After(2 * time.Second):
		c.JSON(http.StatusAccepted, gin.H{
			"message": "File upload is being processed",
			"status":  "processing",
		})
	}
}

// Helper function for concurrent file saving
func saveUploadedFileConcurrently(file *multipart.FileHeader, dst string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	// Use buffered channel for progress tracking if needed
	errChan := make(chan error, 1)

	go func() {
		_, err := io.Copy(out, src)
		errChan <- err
	}()

	return <-errChan
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

	// Check cache first
	cachedFile, err := cache.GetFileMetadata(fileID)
	if err == nil {
		var file database.File
		if err := json.Unmarshal([]byte(cachedFile), &file); err == nil {
			c.JSON(http.StatusOK, file)
			return
		}
	}

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
