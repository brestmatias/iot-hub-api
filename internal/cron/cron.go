package cron

import (
	"iot-hub-api/internal/config"
	"iot-hub-api/internal/station"
	"iot-hub-api/model"
	"log"

	"github.com/robfig/cron/v3"
)

type Cron struct {
	Cron           *cron.Cron
	StationService *station.StationService
	CronService    *CronService
	Config         *config.ConfigFile
	Funcs          []model.CronFuncDTO
}

func New(stationService *station.StationService, cronService *CronService, config *config.ConfigFile) *Cron {
	log.Printf("🎬 🎬 Starting cron")
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
	// TODO !!!! ajustar reloar para que no recargue tareas que no cambian
	// el problema que ocurre es que se ejecuta el reload antes que la tarea corra
	// provoca que nunca se ejecute la tarea

	for _, entry := range s.Cron.Entries() {
		s.Cron.Remove(entry.ID)
	}
	cronTasksToStart := s.CronService.BuildTasks()
	if cronTasksToStart != nil && len(*cronTasksToStart) > 0 {
		for _, i := range *cronTasksToStart {
			log.Printf("⏲️ ⏲️ [cron_task:%v][spec:%v] Cron Task configured.", i.Description, i.Spec)
			newId, _ := s.Cron.AddFunc(i.Spec, i.Func)
			i.EntryID = &newId
			s.Funcs = append(s.Funcs, i)
		}
	}

	//DESABILITO EL AUTORELOAD DEL CRON HASTA QUE SOLUCIONE LO QUE ESCRIBO ARRIBA
	/*if s.Config.Cron.ReloadTaskSpec != "" {
		log.Println("Cron Config to be reloaded ", s.Config.Cron.ReloadTaskSpec)
		s.Cron.AddFunc(s.Config.Cron.ReloadTaskSpec, s.LoadFuncs)
	}*/
}
