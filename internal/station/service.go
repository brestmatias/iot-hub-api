package station

import (
	"fmt"
	"iot-hub-api/internal/network"
	"iot-hub-api/internal/repository"
	"iot-hub-api/internal/restclient"
	"iot-hub-api/model"
	"log"
)

type StationService interface {
	SeekOnlineStations() *[]model.Station
	SeekAndSaveOnlineStations() *[]model.Station
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

func (s *stationService) SeekOnlineStations() *[]model.Station {
	method := "SeekOnlineStations"
	var result []model.Station
	localNetWorkAddresses, _ := network.GetLocalAddresses()
	var scannedNetworks []string
	for _, localNetworkAddress := range *localNetWorkAddresses {
		if networkAlreadyScanned(&scannedNetworks, localNetworkAddress.Net.String()) {
			continue
		}
		log.Printf("[method:%v][interface_name:%v][ip:%v][net_cidr:%v] Looking for online stations", method, localNetworkAddress.Interface.Name, localNetworkAddress.IP, localNetworkAddress.Net.String())
		ips := network.GetAllNetworkIps(&localNetworkAddress)
		for _, ip := range *ips {
			if ip.Equal(localNetworkAddress.IP) {
				continue
			}
			beaconResponse, _ := s.StationClient.GetBeacon(ip.String())
			/*if err != nil {
				log.Println("Error ")
			}*/
			if beaconResponse != nil {
				sta := model.Station{
					ID:      beaconResponse.ID,
					IP:      ip.String(),
					Outputs: beaconResponse.Outputs,
				}
				result = append(result, sta)
				fmt.Println(beaconResponse)
			}
		}
		log.Println("END Looking for online stations")
	}
	return &result
}

func (s *stationService) SeekAndSaveOnlineStations() *[]model.Station {
	foundStations := s.SeekOnlineStations()
	var result []model.Station
	log.Println("Merging stations with database")
	for _, foundSta := range *foundStations {
		dbStation := s.StationRepository.FindByStationID(foundSta.ID)
		fmt.Println("db", dbStation)
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
