package taskExecutor

import (
	"fmt"
	"iot-hub-api/internal/config"
	"iot-hub-api/internal/mqtt"
	"iot-hub-api/model"
	"log"
	"time"
)

type TimerTask struct {
	MqttService *mqtt.MqttService
	task        *model.DispatcherTask
	Config      *config.ConfigFile
}

func newTimerTask(task *model.DispatcherTask, mqttService *mqtt.MqttService, config *config.ConfigFile) *TimerTask {
	return &TimerTask{
		task:        task,
		MqttService: mqttService,
		Config:      config,
	}
}

func (t TimerTask) Execute() {
	log.Printf("[doc_id:%v]Executing Dispatcher TimerTask", t.task.DocId)

	duration, err := time.ParseDuration(t.task.Duration)
	if err != nil {
		log.Printf("[doc_id:%v]Error parsing duration", t.task.DocId)
		return
	}

	time, err := time.Parse("15:04:05", t.task.From)
	if err != nil {
			log.Printf("[doc_id:%v]Error parsing 'from'", t.task.DocId)

		return
	}
	time.Location()
	h,m,s:=time.Clock()

	log.Println(h,m,s, duration)

	//TODO IMPLEMEMTAR TIMER TASK!!!!

	body := model.StationCommandBody{
		Interface: t.task.InterfaceId,
		Value:     1,
		Forced:    false,
	}

	topic := fmt.Sprintf(t.Config.Mqtt.StationCommandTopic, t.task.StationId)
	t.MqttService.SpacedPublishCommand(topic, body)
}
