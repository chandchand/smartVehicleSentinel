package main

import (
	"smartVehicleSentinel/config"
	"smartVehicleSentinel/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	// Inisialisasi koneksi Firebase
	config.InitFirebase()

	r := gin.Default()
	routes.SetupRoutes(r)
	r.Run(":8080")
}
