package routes

import (
	"time"

	"github.com/gin-gonic/gin"
)

func GetModInfoHandler(c *gin.Context) {
	time.Sleep(2 * time.Second)
	c.JSON(200, gin.H{"message": "Modding is not supported"})
}

func RegisterModdingRoutes(group *gin.RouterGroup) {
	group.GET("/GetModInfo", GetModInfoHandler)
}
