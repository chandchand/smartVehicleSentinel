package main

import (
	"context"
	"log"
	"smartVehicleSentinel/config"
	"smartVehicleSentinel/routes"
	"smartVehicleSentinel/services"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gin-gonic/gin"
)

func main() {
	config.InitFirebase()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go services.StartScheduler(ctx)

	opts := config.GetMQTTOptions()
	opts.OnConnect = func(client mqtt.Client) {
		log.Println("✅ MQTT client berhasil terhubung ke broker")

		// Pindahkan pemanggilan SubscribeRelayTopic di sini
		services.SubscribeRelayTopic()

		if token := client.Subscribe("rfid/scan", 0, services.HandleRFIDMessage); token.Wait() && token.Error() != nil {
			log.Println("❌ Gagal subscribe rfid/scan:", token.Error())
		}
	}

	client := mqtt.NewClient(opts)

	// Cek koneksi dan atur config.MQTTClient
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("❌ Gagal koneksi ke MQTT broker: %v", token.Error())
	}

	config.MQTTClient = client

	r := gin.Default()
	routes.SetupRoutes(r)

	log.Println("🚀 Server berjalan di http://localhost:8080")

	// Jalankan server Gin
	r.Run(":8080")
}
