package mqtt

import (
	"iot-hub-api/internal/config"
	"iot-hub-api/internal/repository"
	"log"
)

type MqttController struct {
	MqttService                   *MqttService
	InterfaceLastStatusRepository *repository.InterfaceLastStatusRepository
	Config *config.ConfigFile
}

func NewMqttController(mqttService *MqttService, interfaceLastStatusRepository *repository.InterfaceLastStatusRepository, config *config.ConfigFile) *MqttController {
	i:= &MqttController{
		MqttService:                   mqttService,
		InterfaceLastStatusRepository: interfaceLastStatusRepository,
		Config: config,
	}
	i.subscribe()
	return i
}

func (m *MqttController) subscribe() {
	log.Printf("Subscribing to Station News Topic.")
	if token := m.MqttService.Client.Subscribe(NewStationNewsConsumer(m.InterfaceLastStatusRepository, m.MqttService, m.Config)); token.Wait() && token.Error() != nil {
		log.Println("Error subscribing to station/news", token.Error())
	}
	log.Printf("End subscribing MQTT topics.")

}
