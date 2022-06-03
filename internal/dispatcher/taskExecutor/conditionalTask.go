package taskExecutor

import (
	"iot-hub-api/internal/config"
	"iot-hub-api/internal/mqtt"
	"iot-hub-api/model"
	"log"
)

type ConditionalTask struct {
	MqttService *mqtt.MqttService
	task        *model.DispatcherTask
	Config      *config.ConfigFile
}

func newConditionalTask(task *model.DispatcherTask, mqttService *mqtt.MqttService, config *config.ConfigFile) *ConditionalTask {
	return &ConditionalTask{
		task:        task,
		MqttService: mqttService,
		Config:      config,
	}
}

func (t ConditionalTask) Execute() {
	log.Printf("[doc_id:%v]Executing Dispatcher ConditionalTask", t.task.DocId)
}
