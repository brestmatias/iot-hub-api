package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BeaconResponse struct {
	ID         string   `json:"id"`
	Interfaces []string `json:"interfaces"`
	Broker     string   `json:"broker"`
}

type StationPutResponse struct {
	ID     string `json:"id,omitempty"`
	Broker string `json:"broker,omitempty"`
}

type Station struct {
	DocId               primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	ID                  string             `bson:"id" json:"id"`
	IP                  string             `bson:"ip" json:"ip"`
	Broker              string             `bson:"broker" json:"broker"`
	LastUpdate          primitive.DateTime `bson:"last_update" json:"last_update"`
	Interfaces          []string           `bson:"interfaces" json:"interfaces"`
	LastHandShake       primitive.DateTime `bson:"last_handshake" json:"last_handshake"`
	LastOkHandShake     primitive.DateTime `bson:"last_ok_handshake" json:"last_ok_handshake"`
	LastHandShakeResult string             `bson:"last_handshake_result" json:"last_handshake_result"`
	LastPingStatus      string             `bson:"last_ping_status" json:"last_ping_status"`
}

type InterfaceLastStatus struct {
	DocId           primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	StationID       string             `bson:"station_id" json:"station_id"`
	IntefaceID      string             `bson:"interface_id" json:"interface_id"`
	DispatcherValue int                `bson:"dispatcher_value" json:"dispatcher_value"`
	ReportedValue   int                `bson:"reported_value" json:"reported_value"`
	LastUpdate      primitive.DateTime `bson:"last_update" json:"last_update"`
	LastReport      primitive.DateTime `bson:"last_report" json:"last_report"`
}

type StationCommandBody struct {
	Interface string `json:"interface,omitempty"`
	Value     int    `json:"value,omitempty"`
	Forced    bool   `json:"forced,omitempty"`
}

type StationNewsBody struct {
	Id         string                           `json:"id"`
	Status     string                           `json:"status,omitempty"`
	Interfaces []StationNewsInterfaceStatusBody `json:"interfaces,omitempty"`
}

type StationNewsInterfaceStatusBody struct {
	Id    string `json:"id"`
	Value int    `json:"value,omitempty"`
}
