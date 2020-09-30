package mqttUtils

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"log"
	"mqtt_consumer/config"
	"time"
)

func ConnectToMQTT(cfg config.Mqtt) mqtt.Client {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(cfg.Url)
	opts.SetUsername(cfg.Username)
	opts.SetPassword(cfg.Password)

	client := mqtt.NewClient(opts)
	token := client.Connect()
	for !token.WaitTimeout(3 * time.Second) {
	}
	if err := token.Error(); err != nil {
		log.Fatalf("Failed to connect to MQTT Broker at %v", cfg.Url)
	}
	log.Println("Connected to MQTT Broker...")
	return client
}
