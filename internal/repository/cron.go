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

type CronRepository interface {
	FindAll() *[]model.CronTask
	FindByField(field string, value interface{}) *[]model.CronTask
}

type cronRepository struct {
	MongoDB    *mongo.Database
	Collection *mongo.Collection
}

func (r *cronRepository) FindAll() *[]model.CronTask {
	method := "FindAll"
	cur, err := r.Collection.Find(context.TODO(), bson.D{}, nil)
	if err != nil {
		log.Println(fmt.Errorf("[method:%s]Error Getting Results", method), err)
	}
	defer cur.Close(context.TODO())
	var result []model.CronTask
	err = cur.All(context.TODO(), &result)
	if err != nil {
		log.Println(fmt.Errorf("[method:%s]Error Decoding Results", method), err)
	}
	return &result
}

func (h *cronRepository) FindByField(field string, value interface{}) *[]model.CronTask {
	var result []model.CronTask
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

func NewCronRepository(mongodb *mongo.Database) CronRepository {
	return &cronRepository{
		MongoDB:    mongodb,
		Collection: mongodb.Collection("cron_task"),
	}
}
