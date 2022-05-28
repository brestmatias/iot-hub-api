package cron

import (
	"iot-hub-api/internal/station"
	"log"

	"github.com/robfig/cron/v3"
)

type Cron struct {
	Cron           *cron.Cron
	StationService *station.StationService
	CronService    *CronService
}

func New(stationService *station.StationService, cronService *CronService) *Cron {
	log.Println("Starting Scheduler")
	cron := Cron{
		Cron:           cron.New(),
		StationService: stationService,
		CronService:    cronService,
	}
	cron.Start()
	return &cron
}

func (s *Cron) Start() {
	cronTasksToStart := s.CronService.BuildTasks()
	if cronTasksToStart != nil && len(*cronTasksToStart) > 0 {
		for _, i := range *cronTasksToStart {
			s.Cron.AddFunc(i.Spec, i.Func)
		}
		s.Cron.Start()
	}
}
