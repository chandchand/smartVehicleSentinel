package routes

import (
	"smartVehicleSentinel/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	api := router.Group("/api")
	{
		api.GET("/relay", controllers.GetRelayStatus)
		api.PATCH("/relay", controllers.UpdateRelayStatus)
	}
}
