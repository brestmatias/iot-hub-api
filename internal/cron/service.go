package cron

import (
	cron_tasks "iot-hub-api/internal/cron/tasks"
	"iot-hub-api/internal/repository"
	"iot-hub-api/internal/station"
	"iot-hub-api/model"
	"log"
)

type CronService struct {
	CronRepository *repository.CronRepository
	StationService *station.StationService
}

func NewCronService(cronRepository *repository.CronRepository, stationService *station.StationService) CronService {
	return CronService{
		CronRepository: cronRepository,
		StationService: stationService,
	}
}

func (s CronService) BuildTasks() *[]model.CronFuncDTO {
	tasksDb := (*s.CronRepository).FindByField("enabled", true)
	var result []model.CronFuncDTO

	for _, t := range *tasksDb {
		switch t.TaskId {
		case "seek_stations":
			result = append(result, model.CronFuncDTO{Spec: t.Spec, Func: cron_tasks.NewSeekStationsTask(s.StationService, &t)})
		case "handshake_stations":
			result = append(result, model.CronFuncDTO{Spec: t.Spec, Func: cron_tasks.NewHandshakeTask(s.StationService, &t)})
		case "ping_stations":
			result = append(result, model.CronFuncDTO{Spec: t.Spec, Func: cron_tasks.NewPingTask(s.StationService, &t)})
		default:
			log.Println("Build Cron Task", t.TaskId, "unimplemented!!!!")
		}

	}

	return &result
}
