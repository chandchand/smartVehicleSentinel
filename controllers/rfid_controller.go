package controllers

import (
	"net/http"
	"smartVehicleSentinel/config"
	"smartVehicleSentinel/utils"

	"github.com/gin-gonic/gin"
)

func EnterRFIDRegisterMode(c *gin.Context) {
	var req struct {
		Nama string `json:"nama"`
	}

	if err := c.ShouldBindJSON(&req); err != nil || req.Nama == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nama harus diisi"})
		return
	}

	utils.SetRFIDRegister(req.Nama)

	// Publish ke MQTT untuk masuk mode daftar
	config.MQTTClient.Publish("rfid/mode", 1, false, "register")

	c.JSON(http.StatusOK, gin.H{
		"message": "Berhasil masuk mode daftar RFID",
		"nama":    req.Nama,
	})
}
