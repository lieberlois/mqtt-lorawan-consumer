package main

import (
	"bytes"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"log"
	"mqtt_consumer/config"
	"mqtt_consumer/mqttUtils"
	"mqtt_consumer/parser"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

var cfg config.Config
var writeUrl string

func main() {
	config.LoadConfig(&cfg)

	persistData := len(cfg.InfluxDB.Url) > 0 && len(cfg.InfluxDB.Database) > 0

	if persistData {
		writeUrl = fmt.Sprintf(
			"%s/write?db=%s",
			cfg.InfluxDB.Url,
			cfg.InfluxDB.Database,
		)
	}

	jsonParser := parser.NewParser(cfg.Parser)

	client := mqttUtils.ConnectToMQTT(cfg.MqttBroker)
	token := handleMQTTSubscription(client, jsonParser, persistData)
	token.Wait()
	if err := token.Error(); err != nil {
		log.Fatalf("Could not subscribe to MQTT Topic at %v", cfg.MqttBroker)
	}

	exitSignal := make(chan os.Signal)
	signal.Notify(exitSignal, syscall.SIGINT, syscall.SIGTERM)
	<-exitSignal
}

func handleMQTTSubscription(client mqtt.Client, parser *parser.Parser, persistData bool) mqtt.Token {
	token := client.Subscribe(cfg.MqttBroker.Topic, 0, func(client mqtt.Client, msg mqtt.Message) {
		payload := string(msg.Payload())
		json, err := parser.StringToJson(payload)
		if err != nil {
			log.Printf("Failed to convert %v into a JSON object", payload)
		}
		lineProtocol, err := parser.JsonToInfluxLineProtocol(json)
		if err != nil {
			log.Println(err.Error())
			return
		}
		log.Printf("Received data: %s", lineProtocol)

		if !persistData {
			return
		}
		statusCode, status, err := postDataToInflux(lineProtocol)
		if err != nil {
			log.Printf("Failed to send data %s to the database: %s", lineProtocol, err.Error())
		} else {
			switch statusCode {
			case 204:
				log.Printf("Data written to InfluxDB.")
			default:
				log.Printf("Something went wrong writing to the database: %s", status)
			}
		}
	})
	return token
}

func postDataToInflux(lineProtocol string) (int, string, error) {
	client := &http.Client{}
	buf := new(bytes.Buffer)
	buf.Write([]byte(lineProtocol))

	req, err := http.NewRequest("POST", writeUrl, buf)
	if err != nil {
		return 0, "", err
	}
	req.Header.Add(
		"Authorization",
		fmt.Sprintf(
			"Token %s:%s",
			cfg.InfluxDB.Username,
			cfg.InfluxDB.Password,
		),
	)

	body, err := client.Do(req)

	if err != nil {
		return 0, "", err
	}
	return body.StatusCode, body.Status, nil
}
