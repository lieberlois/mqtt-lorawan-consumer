package mqttUtils

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"log"
	"mqtt_consumer/config"
	"time"
)

func ConnectToMQTT(cfg config.Mqtt, subscribeFunction func(c mqtt.Client)){
	opts := mqtt.NewClientOptions()
	opts.AddBroker(cfg.Url)
	opts.SetUsername(cfg.Username)
	opts.SetPassword(cfg.Password)
	opts.SetCleanSession(false)
	opts.SetClientID(cfg.ClientId)
	opts.OnConnect = subscribeFunction

	client := mqtt.NewClient(opts)
	token := client.Connect()
	for !token.WaitTimeout(10 * time.Second) {
	}
	if err := token.Error(); err != nil {
		log.Fatalf("Failed to connect to MQTT Broker at %v, %s", cfg.Url, err.Error())
	}
}
