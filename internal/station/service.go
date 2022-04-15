package station

import (
	"iot-hub-api/internal/network"
	"iot-hub-api/internal/repository"
	"iot-hub-api/internal/restclient"
	"iot-hub-api/model"
	"iot-hub-api/tracing"
	"net"

	"github.com/gin-gonic/gin"
)

type StationService interface {
	SeekOnlineStations(*gin.Context) *[]model.Station
	SeekAndSaveOnlineStations(*gin.Context) *[]model.Station
}

type stationService struct {
	StationRepository repository.StationRepository
	StationClient     restclient.StationClient
}

func NewStationService(stationRepository repository.StationRepository, stationClient restclient.StationClient) StationService {
	return &stationService{
		StationRepository: stationRepository,
		StationClient:     stationClient,
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
			tracing.Log("[method:%v][interface_name:%v][ip:%v][net_cidr:%v] Network Already scanned", c, method, localNetworkAddress.Interface.Name, localNetworkAddress.IP, localNetworkAddress.Net.String())
			continue
		}
		tracing.Log("[method:%v][interface_name:%v][ip:%v][net_cidr:%v] Looking for online stations", c, method, localNetworkAddress.Interface.Name, localNetworkAddress.IP, localNetworkAddress.Net.String())
		ips := network.GetAllNetworkIps(&localNetworkAddress)
		for _, ip := range *ips {
			if shouldDiscardIp(ip, localNetWorkAddresses) {
				tracing.Log("[method:%v][ip:%v]Discarding IP", c, method, ip)
				continue
			}
			beaconResponse, _ := s.StationClient.GetBeacon(c, ip.String())
			if beaconResponse != nil {
				sta := model.Station{
					ID:      beaconResponse.ID,
					IP:      ip.String(),
					Outputs: beaconResponse.Outputs,
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

//TODO: IMPLEMENTAR UN HANDSHAKE QUE SINCRONICE A TODAS LAS ESTACIONES CON SU PADRE
