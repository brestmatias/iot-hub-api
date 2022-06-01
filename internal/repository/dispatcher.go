package repository

import (
	"context"
	"errors"
	"fmt"
	"iot-hub-api/model"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type DispatcherRepository interface {
	FindAll() *[]model.DispatcherTask
	FindByField(field string, value interface{}) *[]model.DispatcherTask
}

type dispatcherRepository struct {
	MongoDB    *mongo.Database
	Collection *mongo.Collection
}

func (r *dispatcherRepository) FindAll() *[]model.DispatcherTask {
	method := "FindAll"
	cur, err := r.Collection.Find(context.TODO(), bson.D{}, nil)
	if err != nil {
		log.Println(fmt.Errorf("[method:%s]Error Getting Results", method), err)
	}
	defer cur.Close(context.TODO())
	var result []model.DispatcherTask
	err = cur.All(context.TODO(), &result)
	if err != nil {
		log.Println(fmt.Errorf("[method:%s]Error Decoding Results", method), err)
	}
	return &result
}

func (h *dispatcherRepository) FindByField(field string, value interface{}) *[]model.DispatcherTask {
	var result []model.DispatcherTask
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

func NewDispatcherRepository(mongodb *mongo.Database) DispatcherRepository {
	return &dispatcherRepository{
		MongoDB:    mongodb,
		Collection: mongodb.Collection("dispatcher_task"),
	}
}
