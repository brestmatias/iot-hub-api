package repository

import (
	"context"
	"errors"
	"iot-hub-api/model"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type HubConfigRepository interface {
	FindAll() *[]model.HubConfig
	FindByHostName(hostName string) *[]model.HubConfig
}

type hubConfigRepository struct {
	MongoDB             *mongo.Database
	HubConfigCollection *mongo.Collection
}

// FindAll implements HubConfigRepository
func (*hubConfigRepository) FindAll() *[]model.HubConfig {
	panic("unimplemented")
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
	err = findResult.Decode(&result)
	if err != nil {
		log.Println(err)
	}
	return &result
}

func NewHubConfigRepository(mongodb *mongo.Database) HubConfigRepository {
	hubConfigCollection := mongodb.Collection("station")
	return &hubConfigRepository{
		MongoDB:             mongodb,
		HubConfigCollection: hubConfigCollection,
	}
}
