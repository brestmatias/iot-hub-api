package restclient

import (
	"iot-hub-api/model"
	"iot-hub-api/tracing"

	"github.com/brestmatias/golang-restclient/rest"
	"github.com/gin-gonic/gin"
)

type StationClient interface {
	GetBeacon(c *gin.Context, address string) (*model.BeaconResponse, error)
	SetBroker(c *gin.Context, stationIP string, value string) (*model.StationPutResponse, error)
}

type stationClient struct {
	rb     *rest.RequestBuilder
	slowRb *rest.RequestBuilder
}

func NewStationClient(requestBuilder *rest.RequestBuilder, slowRequestBuilder *rest.RequestBuilder) StationClient {
	return &stationClient{
		rb:     requestBuilder,
		slowRb: slowRequestBuilder,
	}
}

// GetBeacon implements StationClient
func (s *stationClient) GetBeacon(c *gin.Context, address string) (*model.BeaconResponse, error) {
	if tracing.VerboseOn(c) {
		defer tracing.Un(tracing.Trace(c, "GetBeacon "+address))
	}
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

func (s *stationClient) SetBroker(c *gin.Context, stationIP string, value string) (*model.StationPutResponse, error) {
	if tracing.VerboseOn(c) {
		defer tracing.Un(tracing.Trace(c, "SetBroker "+value))
	}
	var response model.StationPutResponse
	r := s.slowRb.Put(stationIP+"/station", &model.StationPutResponse{Broker: value})
	if r.Err != nil {
		return nil, r.Err
	}
	err := r.FillUp(&response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}
