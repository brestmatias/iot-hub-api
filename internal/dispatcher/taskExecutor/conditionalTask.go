package taskExecutor

import (
	"iot-hub-api/internal/mqtt"
	"iot-hub-api/model"
	"log"
)

type ConditionalTask struct {
	MqttService *mqtt.MqttService
	task        *model.DispatcherTask
}

func newConditionalTask(task *model.DispatcherTask, mqttService *mqtt.MqttService) *ConditionalTask {
	return &ConditionalTask{
		task:        task,
		MqttService: mqttService,
	}
}

func (t ConditionalTask) Execute() {
	log.Printf("[doc_id:%v]Executing Dispatcher ConditionalTask", t.task.DocId)
}
