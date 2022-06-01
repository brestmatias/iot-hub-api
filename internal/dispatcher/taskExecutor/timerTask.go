package taskExecutor

import (
	"fmt"
	"iot-hub-api/internal/mqtt"
	"iot-hub-api/model"
	"log"
)

type TimerTask struct {
	MqttService *mqtt.MqttService
	task        *model.DispatcherTask
}

func newTimerTask(task *model.DispatcherTask, mqttService *mqtt.MqttService) *TimerTask {
	return &TimerTask{
		task:        task,
		MqttService: mqttService,
	}
}

func (t TimerTask) Execute() {
	log.Printf("[doc_id:%v]Executing Dispatcher TimerTask", t.task.DocId)

	body := model.StationCommandBody{
		Interface: t.task.InterfaceId,
		Value:     1,
		Forced:    false,
	}

	topic := fmt.Sprintf("command/%s", t.task.StationId)
	t.MqttService.PublishCommand(topic, body)
}
