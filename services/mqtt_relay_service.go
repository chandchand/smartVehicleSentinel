package services

import (
	"fmt"
	"log"
	"smartVehicleSentinel/config"
	"smartVehicleSentinel/models"
	"smartVehicleSentinel/utils"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var currentRelayStatus models.Relay // Variabel global untuk menyimpan status
func SubscribeRelayTopic() {
	if config.MQTTClient == nil || !config.MQTTClient.IsConnected() {
		log.Println("❌ MQTT Client belum terhubung, tidak bisa subscribe sekarang")
		return
	}

	token := config.MQTTClient.Subscribe("vehicle/relay", 0, func(client mqtt.Client, msg mqtt.Message) {
		payload := string(msg.Payload())
		log.Printf("📥 Pesan diterima: %s", payload)

		switch payload {
		case "contact_on":
			currentRelayStatus.Contact = true
		case "contact_off":
			currentRelayStatus.Contact = false
		case "engine_on":
			currentRelayStatus.Engine = true
		case "engine_off":
			currentRelayStatus.Engine = false
		case "key_on":
			currentRelayStatus.Key = true
		case "key_off":
			currentRelayStatus.Key = false
		default:
			log.Printf("⚠️ Perintah tidak dikenali: %s", payload)
			return
		}

		// Update Firebase dengan currentRelayStatus terbaru
		relayStatusMap := map[string]interface{}{
			"contact": currentRelayStatus.Contact,
			"engine":  currentRelayStatus.Engine,
			"key":     currentRelayStatus.Key,
			"last_on": time.Now().Format(time.RFC3339), // hasil: 2025-06-11T15:48:21+07:00, // Tambahkan field last_on jika diperlukan
		}

		err := UpdateRelayStatus(relayStatusMap)
		if err != nil {
			log.Printf("❌ Gagal update status relay: %v", err)
		} else {
			log.Printf("✅ Relay status diperbarui: %+v", currentRelayStatus)
		}

	})

	if token.Wait() && token.Error() != nil {
		log.Printf("❌ Gagal subscribe ke topic vehicle/relay : %v ", token.Error())
	} else {
		log.Println("✅ Berhasil subscribe ke topic vehicle/relay")
	}
}

func UpdateRelayStatusFromCommand(command string) error {
	// Ambil status relay sekarang (dari Firebase atau DB)
	currentStatus, err := GetCurrentRelayStatus()
	if err != nil {
		return fmt.Errorf("gagal ambil status relay saat ini: %v", err)
	}

	// Update status sesuai command
	switch command {
	case "contact_on":
		currentStatus.Contact = true
	case "contact_off":
		currentStatus.Contact = false
	case "engine_on":
		currentStatus.Engine = true
	case "engine_off":
		currentStatus.Engine = false
	case "key_on":
		currentStatus.Key = true
	case "key_off":
		currentStatus.Key = false
	default:
		return fmt.Errorf("command tidak dikenali: %s", command)
	}

	currentStatus.LastOnRaw = utils.GetNowInWIB().Format("2006-01-02 15:04:05")
	// currentStatus.LastOnRaw = time.Now().Format(time.RFC3339) // hasil: 2025-06-11T15:48:21+07:00

	relayStatusMap := map[string]interface{}{
		"contact": currentStatus.Contact,
		"engine":  currentStatus.Engine,
		"key":     currentStatus.Key,
		"last_on": currentStatus.LastOnRaw, // Tambahkan field last_on jika diperlukan
	}
	return UpdateRelayStatus(relayStatusMap)
}

func PublishRelayCommand(target string, state string) {
	if config.MQTTClient == nil || !config.MQTTClient.IsConnected() {
		log.Println("❌ MQTT Client belum terhubung")
		return
	}

	command := fmt.Sprintf("%s_%s", target, state) // contoh: "contact_on"
	token := config.MQTTClient.Publish("vehicle/relay", 0, false, command)
	token.Wait()

	if token.Error() != nil {
		log.Printf("❌ Gagal publish ke MQTT: %v", token.Error())
	} else {
		log.Printf("📤 Publish MQTT: %s", command)
		_ = UpdateRelayStatusFromCommand(command) // Update status di Firebase
	}
}
