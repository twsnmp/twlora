package main

import (
	"context"
	"encoding/json"
	"log"
	"strconv"
	"strings"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var mqttCh = make(chan string, 2000)

type mqttPublishDataEnt struct {
	Time               string `json:"time"`
	Type               string `json:"type"`
	ID                 string `json:"id"`
	Detected           bool   `json:"detected"`
	MovingDistance     int    `json:"movingDistance"`
	StationaryDistance int    `json:"stationaryDistance"`
}

func startMQTT(ctx context.Context) {
	if mqttDst == "" {
		return
	}
	log.Println("start mqtt")
	opts := mqtt.NewClientOptions()
	opts.AddBroker(mqttDst)
	if mqttUser != "" && mqttPassword != "" {
		opts.SetUsername(mqttUser)
		opts.SetPassword(mqttPassword)
	}
	opts.SetClientID(mqttClientID)
	opts.SetAutoReconnect(true)
	if debug {
		opts.OnConnect = connectHandler
		opts.OnConnectionLost = connectLostHandler
	}
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Println(token.Error())
		return
	}
	defer client.Disconnect(250)
	for {
		select {
		case <-ctx.Done():
			log.Println("stop mqtt")
			return
		case msg := <-mqttCh:
			if s := makeMqttData(msg); s != "" {
				client.Publish(mqttTopic, 1, false, s).Wait()
			}
		}
	}
}

func makeMqttData(msg string) string {
	// Parse format: RM,ID,status,dist1,dist2 (e.g., RM,1,T,45,20)
	a := strings.Split(msg, ",")
	if len(a) != 5 || a[0] != "RM" {
		if debug {
			log.Printf("skip mqtt message: %s", msg)
		}
		return ""
	}
	var err error
	d := new(mqttPublishDataEnt)
	d.Time = time.Now().Format(time.RFC3339)
	d.Type = a[0]
	d.ID = a[1]
	d.Detected = a[2] == "T"
	if d.MovingDistance, err = strconv.Atoi(a[3]); err != nil {
		return ""
	}
	if d.StationaryDistance, err = strconv.Atoi(a[4]); err != nil {
		return ""
	}
	if j, err := json.Marshal(d); err == nil {
		return string(j)
	}
	return ""
}

func publishMQTT(msg string) {
	select {
	case mqttCh <- msg:
	default:
		if debug {
			log.Println("mqtt channel full, skipping message")
		}
	}
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	log.Println("Connected")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	log.Printf("Connect lost: %v", err)
}
