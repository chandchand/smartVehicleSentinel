package services

import (
	"fmt"
	"log"
	"smartVehicleSentinel/config"
	"smartVehicleSentinel/models"

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
	var status models.Relay

	switch command {
	case "contact_on":
		status.Contact = true
	case "contact_off":
		status.Contact = false
	case "engine_on":
		status.Engine = true
	case "engine_off":
		status.Engine = false
	case "key_on":
		status.Key = true
	case "key_off":
		status.Key = false
	default:
		return fmt.Errorf("command tidak dikenali: %s", command)
	}

	// Kirim ke Firebase atau DB-mu
	statusMap := map[string]interface{}{
		"contact": status.Contact,
		"engine":  status.Engine,
		"key":     status.Key,
	}
	return UpdateRelayStatus(statusMap)
}
