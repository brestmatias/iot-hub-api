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

type StationRepository interface {
	FindByStationID(stationID string) *model.Station
	InsertOne(model.Station) *model.Station
	Update(model.Station) (*model.Station, error)
}

type stationRepository struct {
	MongoDB *mongo.Database
}

func NewStationRepository(mongodb *mongo.Database) StationRepository {
	return &stationRepository{
		MongoDB: mongodb,
	}
}

// InsertOne implements StationRepository
func (s *stationRepository) InsertOne(in model.Station) *model.Station {
	res, err := s.getStationCollection().InsertOne(context.Background(), in)
	if err != nil {
		log.Println(err)
	}
	in.DocId = res.InsertedID.(primitive.ObjectID)

	return &in
}

func (s *stationRepository) Update(in model.Station) (*model.Station, error) {
	collection := s.getStationCollection()
	filter := bson.M{"_id": in.DocId}
	in.LastUpdate = primitive.NewDateTimeFromTime(time.Now())

	update := bson.M{
		"$set": in,
	}
	_, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Println("Error Updating Station", err)
		return nil, err
	}
	return &in, err
}

// FindByStationID implements StationRepository
func (s *stationRepository) FindByStationID(stationID string) *model.Station {
	var result model.Station
	filter := bson.M{"id": stationID}
	findResult := s.getStationCollection().FindOne(context.Background(), filter)
	if findResult.Err() != nil && errors.Is(findResult.Err(), mongo.ErrNoDocuments) {
		return nil
	}
	err := findResult.Decode(&result)
	if err != nil {
		log.Println(err)
	}
	return &result
}

func (s *stationRepository) getStationCollection() *mongo.Collection {
	return s.MongoDB.Collection("station")
}
