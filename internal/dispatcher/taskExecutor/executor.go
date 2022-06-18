package taskExecutor

import (
	"iot-hub-api/internal/config"
	"iot-hub-api/internal/mqtt"
	"iot-hub-api/model"
)

type TaskExecutor interface {
	Execute()
}

func NewExecutor(task *model.DispatcherTask, mqttService *mqtt.MqttService, config *config.ConfigFile, v model.InterfaceLastValueUpdater) TaskExecutor {
	switch task.Type {
	case model.TimerDispatcherTask:
		return newTimerTask(task, mqttService, config, v)
	case model.ConditionalDispatcherTask:
		return newConditionalTask(task, mqttService, config, v)
	default:
		return nil
	}

}
