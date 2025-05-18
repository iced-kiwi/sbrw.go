package routes

import (
	"log"

	"gosbrw/database/structs"
	"gosbrw/services/launcher"

	"github.com/gin-gonic/gin"
)

func UserRegistrationHandler(c *gin.Context) {
	var requestPayload structs.UserRegistrationRequest
	if err := c.ShouldBindJSON(&requestPayload); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request payload: " + err.Error()})
		return
	}
	status, responseBody := launcher.UserRegistration(c, requestPayload)
	c.JSON(status, responseBody)
}

func UserLoginHandler(c *gin.Context) {
	log.Printf("Login request received at: %s %s", c.Request.Method, c.Request.URL.String())
	log.Printf("Content-Type: %s", c.GetHeader("Content-Type"))
	log.Printf("Headers: %v", c.Request.Header)

	var requestPayload structs.UserLoginRequest
	if err := c.ShouldBindJSON(&requestPayload); err != nil {
		log.Printf("Error binding JSON: %v", err)
		c.JSON(400, gin.H{"error": "Invalid request payload: " + err.Error()})
		return
	}
	log.Printf("Login attempt for email: %s", requestPayload.Email)

	status, responseBody := launcher.UserLogin(c, requestPayload)
	log.Printf("Login response: %d - %v", status, responseBody)
	c.JSON(status, responseBody)
}

func RegisterAuthRoutes(group *gin.RouterGroup) {
	group.POST("/modernRegister", UserRegistrationHandler)
	group.POST("/modernAuth", UserLoginHandler)
}
