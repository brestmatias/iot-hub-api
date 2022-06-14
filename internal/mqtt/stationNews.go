package mqtt

import (
	"encoding/json"
	"fmt"
	"iot-hub-api/model"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type StationNewsConsumer struct {
}

func NewStationNewsConsumer() (string, byte, mqtt.MessageHandler) {
	h := &StationNewsConsumer{}
	return "station/news", 0, h.handler
}

func (h StationNewsConsumer) handler(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
	var obj model.StationNewsBody
    if err := json.Unmarshal(msg.Payload(), &obj); err != nil {
        //TODO manejar error
    }
	// TODO manejar novedad en interface
    fmt.Println(obj)
}
