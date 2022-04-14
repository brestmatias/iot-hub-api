package restclient

import (
	"fmt"
	"iot-hub-api/model"

	"github.com/brestmatias/golang-restclient/rest"
)

type StationClient interface {
	GetBeacon(address string) (model.BeaconResponse, error)
}

type stationClient struct {
	rb *rest.RequestBuilder
}

func NewStationClient(requestBuilder *rest.RequestBuilder) StationClient {
	return &stationClient{
		rb: requestBuilder,
	}
}

// GetBeacon implements StationClient
func (s *stationClient) GetBeacon(address string) (model.BeaconResponse, error) {
	var response model.BeaconResponse
	r := s.rb.Get(address + "/beacon")
	fmt.Println(r)
	if r.Err != nil {
		return response, r.Err
	}
	err := r.FillUp(response)
	return response, err
}
