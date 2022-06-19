package mqtt

import (
	"encoding/json"
	"fmt"
	"iot-hub-api/internal/config"
	"iot-hub-api/internal/repository"
	"iot-hub-api/model"
	"log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type StationNewsConsumer struct {
	InterfaceLastStatusRepository *repository.InterfaceLastStatusRepository
	MqttService                   *MqttService
	Config                        *config.ConfigFile
}

func NewStationNewsConsumer(interfaceLastStatusRepository *repository.InterfaceLastStatusRepository, mqttService *MqttService, config *config.ConfigFile) (string, byte, mqtt.MessageHandler) {
	h := &StationNewsConsumer{
		InterfaceLastStatusRepository: interfaceLastStatusRepository,
		MqttService:                   mqttService,
		Config: config,
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
		/**
			TODO!!! Reenviar el estado del dispatcher!!!!
			Fijarse si hay que tener en cuenta el último estado
			Qizás es mejor empezar a enviar el último estado reportado, siempre y cuando sea posterior al estado del dispatcher
			En realidad todo pasa por el dispatcher por lo que el último estado debería coincidir siempre
		*/
		log.Printf("🏁🏁[station_id:%v] Station ready_up message received.",obj.Id)
		result := (*h.InterfaceLastStatusRepository).FindByField("station_id", obj.Id)
		for _, i := range *result {
			topic := fmt.Sprintf(h.Config.Mqtt.StationCommandTopic, i.StationID)
			body := model.StationCommandBody{
				Interface: i.IntefaceID,
				Value:     i.DispatcherValue,
			}
			h.MqttService.PublishCommand(topic, body)
		}
	} else if obj.Status == "interface_update" || obj.Status == "publish_status" {
		for _, i := range obj.Interfaces {
			(*h.InterfaceLastStatusRepository).UpsertReportedStatus(obj.Id, i.Id, i.Value)
		}
	}
}
