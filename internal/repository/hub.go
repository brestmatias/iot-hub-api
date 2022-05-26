package repository

import (
	"context"
	"errors"
	"iot-hub-api/model"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type HubConfigRepository interface {
	FindAll() *[]model.HubConfig
	FindByHostName(hostName string) *[]model.HubConfig
	InsertOne(model.HubConfig) *model.HubConfig
	Update(model.HubConfig) (*model.HubConfig, error)
}

type hubConfigRepository struct {
	MongoDB             *mongo.Database
	HubConfigCollection *mongo.Collection
}

// FindAll implements HubConfigRepository
func (*hubConfigRepository) FindAll() *[]model.HubConfig {
	panic("unimplemented")
}

func (h *hubConfigRepository) InsertOne(config model.HubConfig) *model.HubConfig {
	config.LastUpdate = primitive.NewDateTimeFromTime(time.Now())
	res, err := h.HubConfigCollection.InsertOne(context.Background(), config)

	if err != nil {
		log.Println(err)
	}
	config.DocId = res.InsertedID.(primitive.ObjectID)

	return &config
}

func (h *hubConfigRepository) Update(config model.HubConfig) (*model.HubConfig, error) {
	collection := h.HubConfigCollection
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

// FindByHostName implements HubConfigRepository
func (h *hubConfigRepository) FindByHostName(hostName string) *[]model.HubConfig {
	var result []model.HubConfig
	filter := bson.M{"host_name": hostName}
	findResult, err := h.HubConfigCollection.Find(context.Background(), filter, nil)
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
		HubConfigCollection: hubConfigCollection,
	}
}
