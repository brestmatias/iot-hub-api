package mqtt

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"iot-hub-api/internal/config"
	hub_config "iot-hub-api/internal/hubConfig"
	"iot-hub-api/internal/repository"
	"log"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

type MqttService struct {
	BrokerIp                      string
	Client                        MQTT.Client
	HubConfigService              *hub_config.HubConfigService
	SentCommands                  []CommandHash
	InterfaceLastStatusRepository *repository.InterfaceLastStatusRepository
	Config                        *config.ConfigFile
}

type CommandHash struct {
	Topic    string
	LastHash string
	LastSent time.Time
}

func NewMqttService(hubConfigService *hub_config.HubConfigService, configs *config.ConfigFile, interfaceLastStatusRepository *repository.InterfaceLastStatusRepository) *MqttService {
	service := MqttService{
		HubConfigService:              hubConfigService,
		InterfaceLastStatusRepository: interfaceLastStatusRepository,
		Config:                        configs,
	}
	service.buildClient()
	return &service
}

func (m *MqttService) buildClient() {
	method := "buildClient"

	brokerIp := m.HubConfigService.GetBrokerAddress()
	o := MQTT.NewClientOptions()
	o.AddBroker(fmt.Sprintf("tcp://%v:1883", brokerIp))
	o.SetClientID("iot-dispatcher")
	o.SetUsername("dispatcher")
	o.SetPingTimeout(1 * time.Second)

	m.Client = MQTT.NewClient(o)
	if token := m.Client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("[method:%v] %v", method, token.Error().Error())
	}
}

/*
	Envio de comando manteniendo espacio definido por configuración
	Evita inundar la cola del tópico con envío recurrente del mismo mensaje
	Si el mensaje a enviar es igual al anterior, deberá cumplirse el intervalo de espacio definido
*/
func (m *MqttService) SpacedPublishCommand(topic string, message interface{}) bool {
	method := "SpacedPublishCommand"

	if m.shouldSend(topic, message) == false {
		return false
	}

	messageJSON, _ := json.Marshal(message)
	token := m.Client.Publish(topic, 0, false, messageJSON)
	log.Printf("[method:%v][topic:%v] Command Published", method, topic)
	token.Wait()

	return true
}

func (m *MqttService) shouldSend(topic string, message interface{}) bool {
	h := sha256.New()
	h.Write([]byte(fmt.Sprintf("%v", message)))
	hash := fmt.Sprintf("%x", h.Sum(nil))
	minInterval, _ := time.ParseDuration(m.Config.Mqtt.MinInterval)
	if len(m.SentCommands) == 0 {
		m.SentCommands = append(m.SentCommands, CommandHash{
			Topic:    topic,
			LastHash: hash,
			LastSent: time.Now(),
		})
		return true
	}

	for i, command := range m.SentCommands {
		if command.Topic == topic {
			if command.LastHash == hash {
				if diff := time.Now().Sub(command.LastSent); diff >= minInterval {
					(&m.SentCommands[i]).LastSent = time.Now()
					return true
				} else {
					return false
				}
			} else {
				(&m.SentCommands[i]).LastHash = hash
				(&m.SentCommands[i]).LastSent = time.Now()
				return true
			}
		}
	}
	m.SentCommands = append(m.SentCommands, CommandHash{
		Topic:    topic,
		LastHash: hash,
		LastSent: time.Now(),
	})
	return true
}

func (m *MqttService) PublishCommand(topic string, message interface{}) bool {
	method := "PublishCommand"

	messageJSON, _ := json.Marshal(message)
	if token := m.Client.Connect(); token.Wait() && token.Error() != nil {
		log.Printf("[method:%v] %v", method, token.Error().Error())
		return false
	}
	token := m.Client.Publish(topic, 0, false, messageJSON)
	log.Printf("[method:%v][topic:%v] Command Published", method, topic)
	token.Wait()

	// TODO !!!! OJO con esto que sigue, porque se desconectaba el cliente y se desuscribía
	// revisar si es necesario recheckear la conexión cada x tiempo......
	// sobre todo por las suscripciones
	//m.Client.Disconnect(250) <-- no descomentar, solo lo dejo para acordarme de investigar

	return true
}