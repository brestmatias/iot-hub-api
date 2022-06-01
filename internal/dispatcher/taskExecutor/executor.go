package taskExecutor

import (
	"iot-hub-api/internal/mqtt"
	"iot-hub-api/model"
)

type TaskExecutor interface {
	Execute()
}

func NewExecutor(task *model.DispatcherTask, mqttService *mqtt.MqttService) TaskExecutor {
	switch task.Type {
	case model.TimerDispatcherTask:
		return newTimerTask(task, mqttService)
	case model.ConditionalDispatcherTask:
		return newConditionalTask(task, mqttService)
	default:
		return nil
	}

}
