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

	duration, err := time.ParseDuration(t.task.Duration)
	if err != nil {
		log.Printf("[doc_id:%v]Error parsing duration", t.task.DocId)
		return
	}

	from, err := time.Parse(HMSLayout, t.task.From)
	if err != nil {
		log.Printf("[doc_id:%v]Error parsing 'from'", t.task.DocId)
		return
	}

	//h,m,s:=from.Clock()
	//log.Println(h,m,s, duration.Seconds())

	//TODO IMPLEMEMTAR TIMER TASK!!!!

	if shouldBeOn(from, duration) {
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

func shouldBeOn(from time.Time, duration time.Duration) bool {
	check, _ := time.Parse(HMSLayout, time.Now().Format(HMSLayout))
	return isInTimeSpan(from, duration, check)
}

func isInTimeSpan(from time.Time, duration time.Duration, timeToCheck time.Time) bool {
	end := from.Add(duration)
	log.Println("End:", end)
	return timeToCheck.After(from) && (timeToCheck.Before(end)|| timeToCheck.Equal(end))
}
