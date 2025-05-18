package middleware

import (
	"gosbrw/database"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func Secured() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "No token provided"})
			c.Abort()
			return
		}

		userToken, err := database.GetUserByToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		expiresAt, err := time.Parse(time.RFC3339, userToken.ExpiresAt)
		if err != nil || time.Now().After(expiresAt) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token expired"})
			c.Abort()
			return
		}

		user, err := database.GetUserByID(userToken.UserID)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			c.Abort()
			return
		}

		c.Set("user", user)
		c.Next()
	}
}

func WithSecured(handler gin.HandlerFunc) gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		Secured()(c)
		handler(c)
	})
}
