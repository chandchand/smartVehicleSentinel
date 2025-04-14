package config

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var MQTTClient mqtt.Client

func GetMQTTOptions() *mqtt.ClientOptions {
	broker := "tcp://armadillo-01.rmq.cloudamqp.com:1883"
	clientID := "vehicle-sentinel-backend"
	username := "zdjavzzs:zdjavzzs"
	password := "zjDPz9J5Med1NW6JBLX2zRcKiPQNOL_Q"

	opts := mqtt.NewClientOptions().
		AddBroker(broker).
		SetClientID(clientID).
		SetUsername(username).
		SetPassword(password).
		SetCleanSession(true)

	return opts
}
