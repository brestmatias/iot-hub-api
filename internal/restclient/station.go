package restclient

import (
	"iot-hub-api/model"

	"github.com/brestmatias/golang-restclient/rest"
)

type StationClient interface {
	GetBeacon(address string) (*model.BeaconResponse, error)
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
func (s *stationClient) GetBeacon(address string) (*model.BeaconResponse, error) {
	var response model.BeaconResponse
	r := s.rb.Get(address + "/beacon")
	if r.Err != nil {
		return nil, r.Err
	}
	err := r.FillUp(&response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}
