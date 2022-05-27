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
)

type HubConfigRepository interface {
	FindAll() *[]model.HubConfig
	FindByField(field string, value string) *[]model.HubConfig
	InsertOne(model.HubConfig) *model.HubConfig
	Update(model.HubConfig) (*model.HubConfig, error)
}

type hubConfigRepository struct {
	MongoDB             *mongo.Database
	Collection *mongo.Collection
}

func (r *hubConfigRepository) FindAll() *[]model.HubConfig {
	method := "FindAll"
	cur, err := r.Collection.Find(context.TODO(), bson.D{}, nil)
	if err != nil {
		log.Println(fmt.Errorf("[method:%s]Error Getting Results", method), err)
	}
	defer cur.Close(context.TODO())
	var result []model.HubConfig
	err = cur.All(context.TODO(), &result)
	if err != nil {
		log.Println(fmt.Errorf("[method:%s]Error Decoding Results", method), err)
	}
	return &result
}

func (h *hubConfigRepository) InsertOne(config model.HubConfig) *model.HubConfig {
	config.LastUpdate = primitive.NewDateTimeFromTime(time.Now())
	res, err := h.Collection.InsertOne(context.Background(), config)

	if err != nil {
		log.Println(err)
	}
	config.DocId = res.InsertedID.(primitive.ObjectID)

	return &config
}

func (h *hubConfigRepository) Update(config model.HubConfig) (*model.HubConfig, error) {
	collection := h.Collection
	filter := bson.M{"_id": config.DocId}
	config.LastUpdate = primitive.NewDateTimeFromTime(time.Now())

	update := bson.M{
		"$set": config,
	}
	_, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Println("Error Updating Hub config", err)
		return nil, err
	}
	return &config, err
}

func (h *hubConfigRepository) FindByField(field string, value string) *[]model.HubConfig {
	var result []model.HubConfig
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

func NewHubConfigRepository(mongodb *mongo.Database) HubConfigRepository {
	hubConfigCollection := mongodb.Collection("hub_config")
	return &hubConfigRepository{
		MongoDB:             mongodb,
		Collection: hubConfigCollection,
	}
}
