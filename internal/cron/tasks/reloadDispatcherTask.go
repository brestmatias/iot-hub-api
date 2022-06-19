package cron_tasks

import (
	"iot-hub-api/internal/dispatcher"
	"iot-hub-api/model"
)

type ReloadDispatcherTask struct {
	DispatcherService *dispatcher.DispatcherService
	DBConfig          model.CronTask
}

func NewReloadDispatcherTask(dispatcherService *dispatcher.DispatcherService, config *model.CronTask) func() {
	task := ReloadDispatcherTask{
		DispatcherService: dispatcherService,
		DBConfig:          *config,
	}
	return task.execute
}

func (t *ReloadDispatcherTask) execute() {
	//log.Println("🔃 🔃 Executing Reload Dispatcher Task: ", t.DBConfig.TaskId, "(", t.DBConfig.DocId.String(), ") 🔃 🔃")
	t.DispatcherService.LoadTasks()
}
