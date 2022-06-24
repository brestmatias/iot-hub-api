package station

import (
	hub_config "iot-hub-api/internal/hubConfig"
	"iot-hub-api/internal/network"
	"iot-hub-api/internal/repository"
	"iot-hub-api/internal/restclient"
	"iot-hub-api/model"
	"iot-hub-api/tracing"
	"log"
	"net"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type StationService interface {
	SeekOnlineStations(*gin.Context) *[]model.BeaconResponse
	SeekAndSaveOnlineStations(*gin.Context) *[]model.Station
	DoHandshake(*gin.Context)
	DoPing(*gin.Context)
	GetInterfaceSummary(c *gin.Context) *[]model.InterfaceSummaryResponse
}

type stationService struct {
	StationRepository             *repository.StationRepository
	HubConfigService              *hub_config.HubConfigService
	StationClient                 *restclient.StationClient
	InterfaceLastStatusRepository *repository.InterfaceLastStatusRepository
}

func NewStationService(stationRepository *repository.StationRepository, hubConfigService *hub_config.HubConfigService,
	stationClient *restclient.StationClient, interfaceLastStatusRepository *repository.InterfaceLastStatusRepository) StationService {
	method := "NewStationService"
	log.Printf("[method:%v]🏗️ 🏗️ Building", method)
	return &stationService{
		StationRepository: stationRepository,
		StationClient:     stationClient,
		HubConfigService:  hubConfigService,
		InterfaceLastStatusRepository: interfaceLastStatusRepository,
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

func (s *stationService) SeekOnlineStations(c *gin.Context) *[]model.BeaconResponse {
	method := "SeekOnlineStations"
	var result []model.BeaconResponse
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
			beaconResponse, _ := (*s.StationClient).GetBeacon(c, ip.String())
			if beaconResponse != nil {
				beaconResponse.IP = ip.String()
				result = append(result, *beaconResponse)
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
		dbStation := (*s.StationRepository).FindByField("mac", foundSta.Mac)
		tracing.Log("[method:%v][station_id:%v]Station was found in DB", c, method, foundSta.ID)
		if dbStation != nil {
			dbStation.ID = foundSta.ID
			dbStation.Broker = foundSta.Broker
			dbStation.Interfaces = s.mergeStationInterfaces(*dbStation, foundSta)
			updateResult, _ := (*s.StationRepository).Update(*dbStation)
			result = append(result, *updateResult)
		} else {
			result = append(result, *(*s.StationRepository).InsertOne(foundSta.MapToStation()))
		}
	}
	return &result
}

func (s *stationService) mergeStationInterfaces(dbStation model.Station, beaconResponse model.BeaconResponse) []model.StationInterface {
	merged := dbStation.Interfaces
	for _, r := range beaconResponse.Interfaces {
		notInDb := true
		for _, m := range merged {
			if m.ID == r {
				notInDb = false
				continue
			}
		}
		if notInDb {
			merged = append(merged, model.StationInterface{ID: r})
		}
	}
	return merged
}

func (s *stationService) DoHandshake(c *gin.Context) {
	method := "DoHandshake"
	brokerIp := s.HubConfigService.GetBrokerAddress()
	if brokerIp == "" {
		tracing.Log("[method:%s][result:%+v]No broker Ip is not set", c, method)
		return
	}
	stations := (*s.StationRepository).FindAll()
	for i := range *stations {
		station := (*stations)[i]
		if station.Broker != brokerIp {
			//tracing.Log("[method:%s][station:%+v][ip:%+v]Doing Handshake", c, method, station.ID, station.IP)
			r, err := (*s.StationClient).SetBroker(c, station.IP, brokerIp)
			if err != nil {
				tracing.Log("[method:%s][station:%v][mac:%v]Error Doing handshake %s", c, method, station.ID, station.Mac, err.Error())
				station.LastHandShakeResult = "error"
			} else {
				tracing.Log("[method:%s][result:%+v]Handshake OK", c, method, r)
				station.LastHandShakeResult = "ok"
				station.LastOkHandShake = primitive.NewDateTimeFromTime(time.Now())
			}
			station.LastHandShake = primitive.NewDateTimeFromTime(time.Now())
			(*s.StationRepository).Update(station)
		} else {
			station.LastHandShake = primitive.NewDateTimeFromTime(time.Now())
			station.LastHandShakeResult = "no_action"
			(*s.StationRepository).Update(station)
		}
	}
	tracing.Log("[method:%s]End Handshake", c, method)
}

func (s *stationService) DoPing(c *gin.Context) {
	method := "DoPing"
	stations := (*s.StationRepository).FindAll()
	tracing.Log("[method:%s][stations:%+v]Doing", c, method, stations)
	for _, station := range *stations {
		response := (*s.StationClient).DoPing(c, station.IP)
		if response {
			station.LastPingStatus = "ok"
		} else {
			station.LastPingStatus = "bad"
		}
		(*s.StationRepository).Update(station)
	}
	tracing.Log("[method:%s]End", c, method)
}

func (s *stationService) GetInterfaceSummary(c *gin.Context) *[]model.InterfaceSummaryResponse {
	// method := "GetInterfaceSummary"
	var result []model.InterfaceSummaryResponse
	stations := (*s.StationRepository).FindAll()
 	//lastStatus:=(*s.InterfaceLastStatusRepository).()

	for _,station:=range *stations{
		for _,interf:=range station.Interfaces {
			result = append(result, model.InterfaceSummaryResponse{
				StationID:   station.ID,
				InterfaceID: interf.ID,
				Name:        interf.Name,
				Description: interf.Description,
				Value:       0,
			})
		}
	}
	return &result
}