package cron

import (
	"iot-hub-api/internal/config"
	"iot-hub-api/internal/station"
	"log"

	"github.com/robfig/cron/v3"
)

type Cron struct {
	Cron           *cron.Cron
	StationService *station.StationService
	CronService    *CronService
	Config         *config.ConfigFile
}

func New(stationService *station.StationService, cronService *CronService, config *config.ConfigFile) *Cron {
	log.Println("Starting Scheduler")
	cron := Cron{
		Cron:           cron.New(),
		StationService: stationService,
		CronService:    cronService,
		Config:         config,
	}
	cron.Start()
	return &cron
}

func (s *Cron) Start() {
	s.LoadFuncs()
	s.Cron.Start()
}

func (s *Cron) LoadFuncs() {
	for _, entry := range s.Cron.Entries() {
		s.Cron.Remove(entry.ID)
	}
	cronTasksToStart := s.CronService.BuildTasks()
	if cronTasksToStart != nil && len(*cronTasksToStart) > 0 {
		for _, i := range *cronTasksToStart {
			log.Printf("⏲️ ⏲️ [cron_task:%v][spec:%v] Cron Task configured.", i.Description, i.Spec)
			s.Cron.AddFunc(i.Spec, i.Func)
		}
	}
	if s.Config.Cron.ReloadTaskSpec != "" {
		log.Println("Cron Config to be reloaded ", s.Config.Cron.ReloadTaskSpec)
		s.Cron.AddFunc(s.Config.Cron.ReloadTaskSpec, s.LoadFuncs)
	}
}
