package routes

import (
	"smartVehicleSentinel/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	api := router.Group("/api")
	{
		api.PATCH("/relay", controllers.SendRelayCommandAndUpdateFirebase) // Endpoint untuk update status relay
		api.POST("/rfid/register", controllers.EnterRFIDRegisterMode)      // Endpoint untuk masuk mode daftar RFID
	}
}
