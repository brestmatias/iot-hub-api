package mqtt

import (
	"encoding/json"
	"fmt"
	"iot-hub-api/internal/repository"
	"iot-hub-api/model"
	"log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type StationNewsConsumer struct {
	InterfaceLastStatusRepository *repository.InterfaceLastStatusRepository
}

func NewStationNewsConsumer(interfaceLastStatusRepository *repository.InterfaceLastStatusRepository) (string, byte, mqtt.MessageHandler) {
	h := &StationNewsConsumer{
		InterfaceLastStatusRepository: interfaceLastStatusRepository,
	}
	return "station/news", 0, h.handler
}

func (h StationNewsConsumer) handler(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
	var obj model.StationNewsBody
	if err := json.Unmarshal(msg.Payload(), &obj); err != nil {
		//TODO manejar error
	}

	if obj.Status == "ready_up" {
		// TODO!!! Reenviar el estado del dispatcher!!!!
		result := (*h.InterfaceLastStatusRepository).FindByField("station_id", obj.Id)
		log.Println(result)
		log.Println("-------------------TODO!!! Reenviar el estado del dispatcher!!!!-----------------")
	} else if obj.Status == "interface_update" || obj.Status == "publish_status" {
		for _, i := range obj.Interfaces {
			(*h.InterfaceLastStatusRepository).UpsertReportedStatus(obj.Id, i.Id, i.Value)
		}
	}
}
