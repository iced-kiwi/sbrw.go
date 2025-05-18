package routes

import (
	"gosbrw/server/routes/middleware"

	"github.com/gin-gonic/gin"
)

func GetPermanentSessionHandler(c *gin.Context) {
	
}

func RegisterUserRoutes(group *gin.RouterGroup) {
	group.POST("/GetPermanentSession", middleware.WithSecured(GetPermanentSessionHandler))
}
