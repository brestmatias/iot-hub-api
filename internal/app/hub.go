package app

import (
	"iot-hub-api/internal/network"
	"iot-hub-api/internal/repository"
	"iot-hub-api/model"
	"log"
)

func MapCurrentHostInterfaces(repo *repository.HubConfigRepository) {
	hostname := network.GetHostName()
	log.Println("Hostname:", hostname)
	nets, _ := network.GetLocalAddresses()
	log.Println("Networks:", nets)

	dbConfigs := (*repo).FindByHostName(hostname)
	var merged []model.HubConfig

	for _, net := range *nets {
		indx := findHubConfigIndexByNet(dbConfigs, &net)
		if indx >= 0 {
			cfg := (*dbConfigs)[indx]
			cfg.Ip = net.IP.String()
			
		} else {

		}

	}
}

func findHubConfigIndexByNet(dbConfigs *[]model.HubConfig, net *network.NetworkAddress) int {
	if dbConfigs == nil {
		return -1
	}
	for i, dbCfg := range *dbConfigs {
		if dbCfg.Interface == net.Interface.Name {
			return i
		}
	}
	return -1
}
