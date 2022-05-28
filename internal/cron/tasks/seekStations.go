package cron_tasks

import (
	"iot-hub-api/internal/station"
	"iot-hub-api/model"
	"log"
)

type SeekStationsTask struct {
	StationService station.StationService
	DBConfig       model.CronTask
}

func NewSeekStationsTask(stationService *station.StationService, config *model.CronTask) func() {
	task := SeekStationsTask{
		StationService: *stationService,
		DBConfig:       *config,
	}
	return task.execute
}

func (t *SeekStationsTask) execute() {
	log.Println("⏲️ ⏲️ Executing Cron Task: ", t.DBConfig.TaskId, "(", t.DBConfig.DocId.String(), ") ⏲️ ⏲️")
	t.StationService.SeekAndSaveOnlineStations(nil)
}
