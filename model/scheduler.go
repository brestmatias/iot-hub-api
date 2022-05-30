package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type DispatcherTask struct {
	DocId       primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	StationId   string             `bson:"station_id,omitempty" json:"station_id,omitempty"`
	InterfaceId string             `bson:"interface_id,omitempty" json:"interface_id,omitempty"`
	From        string             `bson:"from,omitempty" json:"from,omitempty"`
	To          string             `bson:"to,omitempty" json:"to,omitempty"`
	Enabled     bool             `bson:"enabled,omitempty" json:"enabled,omitempty"`
	Command 	string	`bson:"command,omitempty" json:"command,omitempty"`
}
