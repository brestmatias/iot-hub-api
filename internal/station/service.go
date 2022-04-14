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

func (s *stationService) SeekOnlineStations() *[]model.Station {
	var result []model.Station
	netAddresses, _ := network.GetLocalAddresses()
	//TODO!!!! revisar que si el array de networks es de longitud mayor a uno y son las mismas redes hay que hacer una sola recorrida
	for _, i := range *netAddresses {
		log.Println(i.Interface.Name, i.IP)
		ips := network.GetAllNetworkIps(&i)
		log.Println("Looking for alive stations")
		for _, ip := range *ips {
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
		log.Println("END Looking for alive stations")
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
