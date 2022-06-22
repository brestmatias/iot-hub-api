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

type StationRepository interface {
	FindAll() *[]model.Station
	FindByField(field string, value interface{}) *model.Station
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
func (s *stationRepository) FindByField(field string, value interface{}) *model.Station {
	var result model.Station
	filter := bson.M{field: value}
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

func (s *stationRepository) FindAll() *[]model.Station {
	method := "FindAll"
	cur, err := s.getStationCollection().Find(context.TODO(), bson.D{}, nil)
	if err != nil {
		log.Println(fmt.Errorf("[method:%s]Error Getting Results", method), err)
	}
	defer cur.Close(context.TODO())
	var result []model.Station
	err = cur.All(context.TODO(), &result)
	if err != nil {
		log.Println(fmt.Errorf("[method:%s]Error Decoding Results", method), err)
	}
	return &result
}
