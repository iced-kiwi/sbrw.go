package routes

import (
	"github.com/gin-gonic/gin"
)

func ConfigureRoutes(router *gin.Engine) {
	engineSvcGroup := router.Group("/Engine.svc")
	{
		engineSvcGroup.GET("/GetServerInformation", GetServerInformationHandler)
		userGroup := engineSvcGroup.Group("/User")
		{
			RegisterAuthRoutes(userGroup)
			RegisterUserRoutes(userGroup)
		}
		moddingGroup := engineSvcGroup.Group("/Modding")
		{
			RegisterModdingRoutes(moddingGroup)
		}
	}
}
