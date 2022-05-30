package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type CronTask struct {
	DocId   primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	TaskId  string             `bson:"task_id" json:"task_id,omitempty"`
	Spec    string             `bson:"spec" json:"spec,omitempty"`
	Enabled bool               `bson:"enabled" json:"enabled,omitempty"`
}

type CronFuncDTO struct {
	Spec string
	Func func()
}
