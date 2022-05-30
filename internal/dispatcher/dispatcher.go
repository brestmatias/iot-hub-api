package dispatcher

import (
	"iot-hub-api/internal/mqtt"
	"iot-hub-api/model"
)

type DispatcherService struct {
	MqttService *mqtt.MqttService
	Tasks       *model.DispatcherTask[]
}

func NewDispatcherService(mqttService *mqtt.MqttService) *DispatcherService {
	return &DispatcherService{
		MqttService: mqttService,
	}
}


