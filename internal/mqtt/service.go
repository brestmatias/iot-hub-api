package mqtt

import (
	"fmt"
	hub_config "iot-hub-api/internal/hubConfig"
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

func(m *MqttService) PublishCommand (uuid string, topic string ) {
	if token := m.Client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	token := m.Client.Publish("commands/STA01010", 0, false, messageJSON)
	token.Wait()
	m.Client.Disconnect(250)
}

