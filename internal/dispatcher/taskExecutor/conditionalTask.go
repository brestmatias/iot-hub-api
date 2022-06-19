package taskExecutor

import (
	"iot-hub-api/internal/config"
	"iot-hub-api/internal/mqtt"
	"iot-hub-api/model"
)

type ConditionalTask struct {
	MqttService               *mqtt.MqttService
	task                      *model.DispatcherTask
	Config                    *config.ConfigFile
	InterfaceLastValueUpdater model.InterfaceLastValueUpdater
}

func newConditionalTask(task *model.DispatcherTask, mqttService *mqtt.MqttService, config *config.ConfigFile, v model.InterfaceLastValueUpdater) *ConditionalTask {
	return &ConditionalTask{
		task:                      task,
		MqttService:               mqttService,
		Config:                    config,
		InterfaceLastValueUpdater: v,
	}
}

func (t ConditionalTask) Execute() {
	//log.Printf("[doc_id:%v]Executing Dispatcher ConditionalTask", t.task.DocId)
}
