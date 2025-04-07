package controllers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"smartVehicleSentinel/models"
	"strings"

	"github.com/gin-gonic/gin"
)

const firebaseURL = "https://smartvehiclesentinel-2ed68-default-rtdb.asia-southeast1.firebasedatabase.app/relay.json"

// GetRelayStatus membaca status relay dari Firebase
func GetRelayStatus(c *gin.Context) {
	resp, err := http.Get(firebaseURL)
	if err != nil {
		log.Println("❌ Error ambil data:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal ambil data Firebase"})
		return
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	var relay models.Relay
	if err := json.Unmarshal(body, &relay); err != nil {
		log.Println("❌ Error parsing:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal parsing data Firebase"})
		return
	}

	c.JSON(http.StatusOK, relay)
}

// UpdateRelayStatus memperbarui data relay
func UpdateRelayStatus(c *gin.Context) {
	var input models.Relay
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Input tidak valid"})
		return
	}

	payload, _ := json.Marshal(input)
	req, err := http.NewRequest(http.MethodPatch, firebaseURL, strings.NewReader(string(payload)))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat request"})
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengirim request ke Firebase"})
		return
	}
	defer resp.Body.Close()

	c.JSON(http.StatusOK, gin.H{"message": "Berhasil update data"})
}
