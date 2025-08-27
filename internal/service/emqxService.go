package service

import (
	"encoding/json"
	"fmt"

	"github.com/Eatriceeveryday/data-pool-service/internal/entities"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type EmqxService struct {
	client mqtt.Client
	rs     *SensorService
}

func NewEmqxService(client mqtt.Client, rs *SensorService) *EmqxService {
	return &EmqxService{client, rs}
}

func (s *EmqxService) Subscribe(topic string) {
	go func() {
		messageHandler := func(client mqtt.Client, msg mqtt.Message) {
			var m entities.Message
			if err := json.Unmarshal(msg.Payload(), &m); err != nil {
				fmt.Println("Error decoding message:", err)
				return
			}
			fmt.Printf("Received: %+v\n", m)

			if err := s.rs.CreateReport(m); err != nil {
				fmt.Println("Error saving report:", err)
			}
			fmt.Printf("Received message on topic [%s]: %s\n", msg.Topic(), string(msg.Payload()))
		}

		if token := s.client.Subscribe(topic, 1, messageHandler); token.Wait() && token.Error() != nil {
			fmt.Println("Subscribe error:", token.Error())
			return
		}
		fmt.Println("Subscribed to topic:", topic)
	}()
}
