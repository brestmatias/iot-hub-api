package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BeaconResponse struct {
	ID         string   `json:"id"`
	Interfaces []string `json:"interfaces"`
	Broker     string   `json:"broker"`
	Mac        string   `json:"mac"`
	IP         string   `json:"ip"`
}

func (b *BeaconResponse) MapToStation() Station {
	var interfaces []StationInterface
	for _, i := range b.Interfaces {
		interfaces = append(interfaces, StationInterface{ID: i})
	}
	sta := Station{
		ID:         b.ID,
		IP:         b.IP,
		Mac:        b.Mac,
		Broker:     b.Broker,
		Interfaces: interfaces,
	}
	return sta
}

type StationPutResponse struct {
	ID     string `json:"id,omitempty"`
	Broker string `json:"broker,omitempty"`
}

type Station struct {
	DocId               primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	ID                  string             `bson:"id" json:"id"`
	IP                  string             `bson:"ip" json:"ip"`
	Mac                 string             `bson:"mac" json:"mac"`
	Broker              string             `bson:"broker" json:"broker"`
	LastUpdate          primitive.DateTime `bson:"last_update" json:"last_update"`
	Interfaces          []StationInterface `bson:"interfaces" json:"interfaces"`
	LastHandShake       primitive.DateTime `bson:"last_handshake" json:"last_handshake"`
	LastOkHandShake     primitive.DateTime `bson:"last_ok_handshake" json:"last_ok_handshake"`
	LastHandShakeResult string             `bson:"last_handshake_result" json:"last_handshake_result"`
	LastPingStatus      string             `bson:"last_ping_status" json:"last_ping_status"`
}

type StationInterface struct {
	ID          string `bson:"id" json:"id"`
	Name        string `bson:"name" json:"name"`
	Description string `bson:"description" json:"description"`
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

type InterfaceSummaryResponse struct {
	StationID   string `bson:"station_id" json:"station_id"`
	InterfaceID string `bson:"interface_id" json:"interface_id"`
	Name        string `bson:"name" json:"name"`
	Description string `bson:"description" json:"description"`
	Value       int    `bson:"value" json:"value"`
}
