package cron_tasks

import (
	"iot-hub-api/internal/station"
	"iot-hub-api/model"
	"log"
)

type HandshakeTask struct {
	StationService station.StationService
	DBConfig       model.CronTask
}

func NewHandshakeTask(stationService *station.StationService, config *model.CronTask) func() {
	task := HandshakeTask{
		StationService: *stationService,
		DBConfig:       *config,
	}
	return task.execute
}

func (t *HandshakeTask) execute() {
	log.Println("⏲️ ⏲️ Executing Cron Task: ", t.DBConfig.TaskId, "(", t.DBConfig.DocId.String(), ") ⏲️ ⏲️")
	t.StationService.DoHandshake(nil)
}
