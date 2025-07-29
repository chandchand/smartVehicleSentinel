package controllers

import (
	"log"
	"net/http"
	"smartVehicleSentinel/config"
	"smartVehicleSentinel/services"

	"github.com/gin-gonic/gin"
)

func SendRelayCommandAndUpdateFirebase(c *gin.Context) {
	var payload struct {
		Command string `json:"command"`
	}

	if err := c.ShouldBindJSON(&payload); err != nil || payload.Command == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or missing command"})
		return
	}

	command := payload.Command

	// 🔌 Publish ke MQTT
	token := config.MQTTClient.Publish("vehicle/relay", 0, false, command)
	token.Wait()
	if token.Error() != nil {
		log.Printf("❌ Gagal publish ke MQTT: %v", token.Error())
	} else {
		log.Printf("✅ Perintah berhasil dikirim ke MQTT: %s", command)
	}

	// 🔄 Update Firebase sesuai command yang dikirim
	err := services.UpdateRelayStatusFromCommand(command)
	if err != nil {
		log.Printf("❌ Gagal update Firebase: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal update ke Firebase"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Command sent and Firebase updated",
		"command": command,
	})
}

// GetRelayStatus membaca status relay dari Firebase
// func GetRelayStatus(c *gin.Context) {
// 	relay, err := services.GetRelayStatus()
// 	if err != nil {
// 		log.Println("❌ Error ambil data:", err)
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal ambil data Firebase"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, relay)
// }

// UpdateRelayStatus memperbarui data relay
// func UpdateRelayStatus(c *gin.Context) {
// 	var input models.Relay
// 	if err := c.ShouldBindJSON(&input); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Input tidak valid"})
// 		return
// 	}

// 	err := services.UpdateRelayStatus(input)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengirim request ke Firebase"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"message": "Berhasil update data"})
// }
