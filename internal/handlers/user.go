package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/manojkp08/22BCE11415_Backend/internal/auth"
	"github.com/manojkp08/22BCE11415_Backend/internal/database"
)

func GoogleLoginHandler(c *gin.Context) {
	auth.HandleGoogleLogin(c.Writer, c.Request)
}

func GoogleCallbackHandler(c *gin.Context) {
	userInfo, err := auth.HandleGoogleCallback(c.Writer, c.Request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Check if user exists in DB, if not create
	user, err := database.GetOrCreateUser(userInfo.Email, userInfo.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// Generating JWT token
	token, err := auth.GenerateJWTToken(user.ID, user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user":  user,
	})
}
