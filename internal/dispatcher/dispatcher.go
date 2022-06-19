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
	MqttService                   *mqtt.MqttService
	DispatcherRepository          *repository.DispatcherRepository
	InterfaceLastStatusRepository *repository.InterfaceLastStatusRepository
	Tasks                         *[]model.DispatcherTask
	Config                        *config.ConfigFile
}

func NewDispatcherService(mqttService *mqtt.MqttService, dispatcherRepository *repository.DispatcherRepository, interfaceLastStatusRepository *repository.InterfaceLastStatusRepository, config *config.ConfigFile) *DispatcherService {
	return &DispatcherService{
		MqttService:                   mqttService,
		DispatcherRepository:          dispatcherRepository,
		InterfaceLastStatusRepository: interfaceLastStatusRepository,
		Config:                        config,
	}
}

/*
	!!!TODO: implementar reload cada x cantidad de segundos, quizás haya que implementar mutex
*/
func (d *DispatcherService) LoadTasks() {
	// TODO hacer que arme un mapa [taskType, []tasks]
	d.Tasks = (*d.DispatcherRepository).FindByField("enabled", true)
	for _, i := range *d.Tasks {
		log.Printf("👷👷[task_type:%s][from:%v][duration:%v]Dispatcher task configured\n", i.Type, i.From, i.Duration)
	}
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
	executor := taskExecutor.NewExecutor(&task, d.MqttService, d.Config, d.updateIntefaceLastStatus)
	if executor == nil {
		log.Println("Executor for ", task.Type, " not implemented!!!!!")
		return
	}
	executor.Execute()
}

func (d DispatcherService) updateIntefaceLastStatus(stationId string, interfaceId string, value int) {
	//TODO ver de manejar una cache para no estar enviando al repo al pedo
	(*d.InterfaceLastStatusRepository).UpsertDispatcherStatus(stationId, interfaceId, value)
}
