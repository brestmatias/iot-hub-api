package station

import (
	"iot-hub-api/internal/network"
	"iot-hub-api/internal/repository"
	"iot-hub-api/internal/restclient"
	"iot-hub-api/model"
	"iot-hub-api/tracing"
	"net"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type StationService interface {
	SeekOnlineStations(*gin.Context) *[]model.Station
	SeekAndSaveOnlineStations(*gin.Context) *[]model.Station
	DoHandshake(*gin.Context)
	DoPing(*gin.Context)
}

type stationService struct {
	StationRepository   repository.StationRepository
	HubConfigRepository repository.HubConfigRepository
	StationClient       restclient.StationClient
}

func NewStationService(stationRepository repository.StationRepository, hubConfigRepository repository.HubConfigRepository, stationClient restclient.StationClient) StationService {
	return &stationService{
		StationRepository:   stationRepository,
		StationClient:       stationClient,
		HubConfigRepository: hubConfigRepository,
	}
}

func networkAlreadyScanned(scannedNetworks *[]string, currentCIDR string) bool {
	if len(*scannedNetworks) == 0 {
		*scannedNetworks = append(*scannedNetworks, currentCIDR)
		return false
	}

	for _, i := range *scannedNetworks {
		if i == currentCIDR {
			return true
		}
	}
	*scannedNetworks = append(*scannedNetworks, currentCIDR)
	return false
}

func shouldDiscardIp(ip net.IP, localAddresses *[]network.NetworkAddress) bool {
	for _, i := range *localAddresses {
		if ip.Equal(i.IP) {
			return true
		}
	}
	return false
}

func (s *stationService) SeekOnlineStations(c *gin.Context) *[]model.Station {
	method := "SeekOnlineStations"
	var result []model.Station
	localNetWorkAddresses, _ := network.GetLocalAddresses()
	var scannedNetworks []string
	for _, localNetworkAddress := range *localNetWorkAddresses {
		if networkAlreadyScanned(&scannedNetworks, localNetworkAddress.Net.String()) {
			// tracing.Log("[method:%v][interface_name:%v][ip:%v][net_cidr:%v] Network Already scanned", c, method, localNetworkAddress.Interface.Name, localNetworkAddress.IP, localNetworkAddress.Net.String())
			continue
		}
		tracing.Log("[method:%v][interface_name:%v][ip:%v][net_cidr:%v] Looking for online stations", c, method, localNetworkAddress.Interface.Name, localNetworkAddress.IP, localNetworkAddress.Net.String())
		ips := network.GetAllNetworkIps(&localNetworkAddress)
		for _, ip := range *ips {
			if shouldDiscardIp(ip, localNetWorkAddresses) {
				// tracing.Log("[method:%v][ip:%v]Discarding IP", c, method, ip)
				continue
			}
			beaconResponse, _ := s.StationClient.GetBeacon(c, ip.String())
			if beaconResponse != nil {
				sta := model.Station{
					ID:         beaconResponse.ID,
					IP:         ip.String(),
					Broker:     beaconResponse.Broker,
					Interfaces: beaconResponse.Interfaces,
				}
				result = append(result, sta)
				tracing.Log("[method:%v][beacon:%+v] Beacon response", c, method, result)
			}
		}
		tracing.Log("[method:%v] END Looking for online stations", c, method)
	}
	return &result
}

func (s *stationService) SeekAndSaveOnlineStations(c *gin.Context) *[]model.Station {
	method := "SeekAndSaveOnlineStations"
	foundStations := s.SeekOnlineStations(c)
	var result []model.Station
	tracing.Log("[method:%v]Merging stations with database", c, method)
	for _, foundSta := range *foundStations {
		dbStation := s.StationRepository.FindByStationID(foundSta.ID)
		tracing.Log("[method:%v][station_id:%v]Station was found in DB", c, method, foundSta.ID)
		if dbStation != nil {
			foundSta.DocId = dbStation.DocId
			updateResult, _ := s.StationRepository.Update(foundSta)
			result = append(result, *updateResult)
		} else {
			result = append(result, *s.StationRepository.InsertOne(foundSta))
		}
	}
	return &result
}

func (s *stationService) DoHandshake(c *gin.Context) {
	method := "DoHandshake"
	brokerIp := s.GetBrokerAddress()
	if brokerIp == "" {
		tracing.Log("[method:%s][result:%+v]No broker Ip is not set", c, method)
		return
	}
	stations := s.StationRepository.FindAll()
	for _, station := range *stations {
		if station.Broker != brokerIp {
			tracing.Log("[method:%s][station:%+v][ip:%+v]Doing Handshake", c, method, station.ID, station.IP)
			r, err := s.StationClient.SetBroker(c, station.IP, brokerIp)
			if err != nil {
				tracing.Log("[method:%s][station:%+v]Error Doing handshake %s", c, method, station, err.Error())
				station.LastHandShakeResult = "error"
			} else {
				tracing.Log("[method:%s][result:%+v]Handshake OK", c, method, r)
				station.LastHandShakeResult = "ok"
				station.LastOkHandShake = primitive.NewDateTimeFromTime(time.Now())
				station.LastPingStatus = "ok"
			}
			station.LastHandShake = primitive.NewDateTimeFromTime(time.Now())
			s.StationRepository.Update(station)
		}
	}
	tracing.Log("[method:%s]End Handshake", c, method)
}

func (s *stationService) DoPing(c *gin.Context) {
	method := "DoPing"
	stations := s.StationRepository.FindAll()
	tracing.Log("[method:%s][stations:%+v]Doing", c, method, stations)
	for _, station := range *stations {
		beaconResponse, _ := s.StationClient.GetBeacon(c, station.IP)
		if beaconResponse != nil {
			station.LastPingStatus = "ok"
		} else {
			station.LastPingStatus = "bad"
		}
		s.StationRepository.Update(station)
	}
	tracing.Log("[method:%s]End", c, method)
}

func (s *stationService) GetBrokerAddress() string {

	configs := s.HubConfigRepository.FindAll()
	for _, config := range *configs {
		if config.IsMQBroker {
			return config.Ip
		}
	}
	return ""
	/*
		log.Printf("[method:%v] Broker configuration not found in repository ","GetBrokerAddress")
		localNetWorkAddresses, _ := network.GetLocalAddresses()
		localIp := (*localNetWorkAddresses)[0].IP
		return localIp.String()*/
}
