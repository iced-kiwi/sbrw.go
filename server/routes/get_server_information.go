package routes

import (
	"database/sql"
	"gosbrw/database"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetServerInformationHandler(c *gin.Context) {
	serverInfo, err := database.GetServerInformation()
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("GetServerInformationHandler: No server information found in database.")
			c.JSON(http.StatusNotFound, gin.H{"error": "Server information not found."})
			return
		}
		log.Printf("GetServerInformationHandler: Error fetching server information: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error."})
		return
	}

	c.JSON(http.StatusOK, serverInfo)
}

func GetServerInformationRoute(router *gin.Engine) {
	router.GET("/server", GetServerInformationHandler)
}
