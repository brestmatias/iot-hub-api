package taskExecutor

import (
	"fmt"
	"iot-hub-api/internal/config"
	"iot-hub-api/internal/mqtt"
	"iot-hub-api/model"
	"log"
	"time"
)

const HMSLayout = "15:04:05"

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

	if t.shouldBeOn() {
		log.Println("ONNNNN")
	} else {
		log.Println("OFFFFFF")
	}

	body := model.StationCommandBody{
		Interface: t.task.InterfaceId,
		Value:     1,
		Forced:    false,
	}

	topic := fmt.Sprintf(t.Config.Mqtt.StationCommandTopic, t.task.StationId)
	t.MqttService.SpacedPublishCommand(topic, body)
}

func (t TimerTask) shouldBeOn() bool {
	duration, err := time.ParseDuration(t.task.Duration)
	if err != nil {
		log.Printf("[doc_id:%v]Error parsing duration", t.task.DocId)
		return false
	}

	from, err := time.Parse(HMSLayout, t.task.From)
	if err != nil {
		log.Printf("[doc_id:%v]Error parsing 'from'", t.task.DocId)
		return false
	}

	check, _ := time.Parse(HMSLayout, time.Now().Format(HMSLayout))
	to, _ := time.Parse(HMSLayout, from.Add(duration).Format(HMSLayout))

	return isInTimeSpan(from, to, check)
}

func isInTimeSpan(from time.Time, to time.Time, check time.Time) bool {
	return (!from.After(to) && (!from.After(check) && !to.Before(check))) ||
		(from.After(to) && !(from.After(check) && to.Before(check)))
}
