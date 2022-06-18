package repository

import (
	"context"
	"errors"
	"fmt"
	"iot-hub-api/model"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type InterfaceLastStatusRepository interface {
	UpsertDispatcherStatus(stationId string, interfaceId string, value int) int64
	UpsertReportedStatus(stationId string, interfaceId string, value int) int64
	 FindByField(field string, value interface{}) *[]model.InterfaceLastStatus
}

type interfaceLastStatusRepository struct {
	MongoDB    *mongo.Database
	Collection *mongo.Collection
}

func NewInterfaceLastStatusRepository(mongodb *mongo.Database) InterfaceLastStatusRepository {
	return &interfaceLastStatusRepository{
		MongoDB:    mongodb,
		Collection: mongodb.Collection("interface_last_status"),
	}
}

func (h *interfaceLastStatusRepository) UpsertDispatcherStatus(stationId string, interfaceId string, value int) int64 {
	method := "UpsertDispatcherStatus"
	filter := bson.M{"interface_id": interfaceId, "station_id": stationId}

	update := bson.D{{"$set", bson.D{
		{"interface_id", interfaceId},
		{"station_id", stationId},
		{"dispatcher_value", value}, 
		{"last_update", primitive.NewDateTimeFromTime(time.Now())},
	}}}
	opts := options.Update().SetUpsert(true)
	result, err := h.Collection.UpdateOne(context.TODO(), filter, update, opts)
	if err != nil {
		log.Println(fmt.Errorf("[method:%s]Error!!", method), err)
	}
	return result.UpsertedCount
}

func (h *interfaceLastStatusRepository) UpsertReportedStatus(stationId string, interfaceId string, value int) int64 {
	method := "UpsertDispatcherStatus"
	filter := bson.M{"interface_id": interfaceId, "station_id": stationId}

	update := bson.D{{"$set", bson.D{
		{"interface_id", interfaceId},
		{"station_id", stationId},
		{"reported_value", value}, 
		{"last_report", primitive.NewDateTimeFromTime(time.Now())},
	}}}
	opts := options.Update().SetUpsert(true)
	result, err := h.Collection.UpdateOne(context.TODO(), filter, update, opts)
	if err != nil {
		log.Println(fmt.Errorf("[method:%s]Error!!", method), err)
	}
	return result.UpsertedCount
}

func (h *interfaceLastStatusRepository) FindByField(field string, value interface{}) *[]model.InterfaceLastStatus{
	var result []model.InterfaceLastStatus
	filter := bson.M{field: value}
	findResult, err := h.Collection.Find(context.Background(), filter, nil)
	if err != nil {
		log.Println(err)
	}
	if findResult.Err() != nil && errors.Is(findResult.Err(), mongo.ErrNoDocuments) {
		return nil
	}

	defer findResult.Close(context.TODO())
	err = findResult.All(context.TODO(), &result)

	if err != nil {
		log.Println(err)
	}
	return &result
}
