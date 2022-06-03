package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DispatcherTaskType string

const (
	TimerDispatcherTask       DispatcherTaskType = "timer"
	ConditionalDispatcherTask                    = "conditional"
)

var DispatcherTaskTypes = []DispatcherTaskType{
	TimerDispatcherTask,
	ConditionalDispatcherTask,
}

type DispatcherTask struct {
	DocId       primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Type        DispatcherTaskType `bson:"type,omitempty" json:"type,omitempty"` // DispatcherTaskType
	StationId   string             `bson:"station_id,omitempty" json:"station_id,omitempty"`
	InterfaceId string             `bson:"interface_id,omitempty" json:"interface_id,omitempty"`
	From        string 			   `bson:"from,omitempty" json:"from,omitempty"` // HH:MM
	Duration    string    		   `bson:"duration,omitempty" json:"duration,omitempty"` //0h5m0s
	Enabled     bool               `bson:"enabled,omitempty" json:"enabled,omitempty"`
}
