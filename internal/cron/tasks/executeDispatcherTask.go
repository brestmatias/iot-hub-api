package cron_tasks

import (
	"iot-hub-api/internal/dispatcher"
	"iot-hub-api/model"
	"log"
)

type ExecuteDispatcherTask struct {
	DispatcherService *dispatcher.DispatcherService
	DBConfig          model.CronTask
}

func NewExecuteDispatcherTask(dispatcherService *dispatcher.DispatcherService, config *model.CronTask) func() {
	task := ExecuteDispatcherTask{
		DispatcherService: dispatcherService,
		DBConfig:          *config,
	}
	return task.execute
}

func (t *ExecuteDispatcherTask) execute() {
	log.Println("⏲️ ⏲️ Executing Cron Task: ", t.DBConfig.TaskId, "(", t.DBConfig.DocId.String(), ") ⏲️ ⏲️")
	for _, i := range model.DispatcherTaskTypes {
		t.DispatcherService.Execute(i)
	}

}
