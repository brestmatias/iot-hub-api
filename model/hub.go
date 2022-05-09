package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type HubConfig struct {
	DocId      primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	HostName   string             `bson:"host_name,omitempty" json:"host_name,omitempty"`
	Interface  string             `bson:"interface,omitempty" json:"interface,omitempty"`
	Ip         string             `bson:"ip,omitempty" json:"ip,omitempty"`
	IsMQBroker bool               `bson:"is_broker,omitempty" json:"is_broker,omitempty"`
}
