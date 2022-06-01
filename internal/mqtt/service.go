package mqtt

import (
	"encoding/json"
	"fmt"
	hub_config "iot-hub-api/internal/hubConfig"
	"log"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

type MqttService struct {
	BrokerIp         string
	Client           MQTT.Client
	HubConfigService *hub_config.HubConfigService
}

func NewMqttService(hubConfigService *hub_config.HubConfigService) *MqttService {
	service := MqttService{
		HubConfigService: hubConfigService,
	}
	service.buildClient()
	return &service
}

func (m *MqttService) buildClient() {
	brokerIp := m.HubConfigService.GetBrokerAddress()
	o := MQTT.NewClientOptions()
	o.AddBroker(fmt.Sprintf("tcp://%v:1883", brokerIp))
	o.SetClientID("iot-dispatcher")
	o.SetUsername("dispatcher")
	o.SetPingTimeout(1 * time.Second)

	m.Client = MQTT.NewClient(o)
}

func (m *MqttService) PublishCommand(topic string, message interface{}) bool {
	method:="PublishCommand"
	if token := m.Client.Connect(); token.Wait() && token.Error() != nil {
		log.Printf("[method:{%v}] %v",method,token.Error().Error())
		return false
	}

	messageJSON, _ := json.Marshal(message)
	// TODO!!! implementar un hash o algun mecanismo para evitar inundar con el mismo req en un intervalo de tiempo
	token := m.Client.Publish(topic, 0, false, messageJSON)
	token.Wait()
	m.Client.Disconnect(250)

	return true
}
