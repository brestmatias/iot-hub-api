package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BeaconResponse struct {
	ID         string   `json:"id"`
	Interfaces []string `json:"interfaces"`
}

type StationPutResponse struct {
	ID     string `json:"id,omitempty"`
	Broker string `json:"broker,omitempty"`
}

type Station struct {
	DocId      primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	ID         string             `bson:"id" json:"id"`
	IP         string             `bson:"ip" json:"ip"`
	LastUpdate primitive.DateTime `bson:"last_update" json:"last_update"`
	Interfaces []string           `bson:"interfaces" json:"interfaces"`
}
