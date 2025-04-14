package utils

import (
	"smartVehicleSentinel/config"
)

func PublishMQTT(topic, message string) {
	token := config.MQTTClient.Publish(topic, 0, false, message)
	token.Wait()
}
