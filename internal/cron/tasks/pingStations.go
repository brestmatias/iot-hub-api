package cron_tasks

import (
	"iot-hub-api/internal/station"
	"iot-hub-api/model"
)

type PingTask struct {
	StationService station.StationService
	DBConfig       model.CronTask
}

func NewPingTask(stationService *station.StationService, config *model.CronTask) func() {
	task := PingTask{
		StationService: *stationService,
		DBConfig:       *config,
	}
	return task.execute
}

func (t *PingTask) execute() {
	//log.Println("⏲️ ⏲️ Executing Cron Task: ", t.DBConfig.TaskId, "(", t.DBConfig.DocId.String(), ") ⏲️ ⏲️")
	t.StationService.DoPing(nil)
}
