package dispatcher

import (
	"iot-hub-api/internal/config"
	"iot-hub-api/internal/dispatcher/taskExecutor"
	"iot-hub-api/internal/mqtt"
	"iot-hub-api/internal/repository"
	"iot-hub-api/model"
	"log"
)

type DispatcherService struct {
	MqttService          *mqtt.MqttService
	DispatcherRepository *repository.DispatcherRepository
	Tasks                *[]model.DispatcherTask
	Config               *config.ConfigFile
}

func NewDispatcherService(mqttService *mqtt.MqttService, dispatcherRepository *repository.DispatcherRepository, config *config.ConfigFile) *DispatcherService {
	return &DispatcherService{
		MqttService:          mqttService,
		DispatcherRepository: dispatcherRepository,
		Config:               config,
	}
}

/*
	!!!TODO: implementar reload cada x cantidad de segundos, quizás haya que implementar mutex
*/
func (d *DispatcherService) LoadTasks() {
	// TODO hacer que arme un mapa [taskType, []tasks]
	d.Tasks = (*d.DispatcherRepository).FindByField("enabled", true)
}

/*
	Ejecuta DispatcherTasks enabled de acuerdo a los types sumistrados
*/
func (d DispatcherService) Execute(taskType model.DispatcherTaskType) {
	for _, i := range *d.Tasks {
		// TODO hacer que LoadTasks arme un mapa [taskType, []tasks] y buscar por key ¿?
		if i.Type == taskType {
			d.executeTask(i)
		}
	}
}

func (d DispatcherService) executeTask(task model.DispatcherTask) {
	executor := taskExecutor.NewExecutor(&task, d.MqttService, d.Config)
	if executor == nil {
		log.Println("Executor for ", task.Type, " not implemented!!!!!")
		return
	}
	executor.Execute()
}
