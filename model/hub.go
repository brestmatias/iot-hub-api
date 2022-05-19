package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type HubConfig struct {
	DocId      primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	HostName   string             `bson:"host_name" json:"host_name,omitempty"`
	Interface  string             `bson:"interface" json:"interface,omitempty"`
	Ip         string             `bson:"ip" json:"ip,omitempty"`
	IsMQBroker bool               `bson:"is_broker" json:"is_broker,omitempty"`
	LastUpdate primitive.DateTime `bson:"last_update" json:"last_update"`
}
