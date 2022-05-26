package app

import (
	"log"

	"iot-hub-api/internal/network"
	"iot-hub-api/internal/repository"
	"iot-hub-api/model"
)

func MapCurrentHostInterfaces(repo repository.HubConfigRepository) {
	hostname := network.GetHostName()
	log.Println("Hostname:", hostname)
	nets, _ := network.GetLocalAddresses()
	log.Println("Networks:", nets)

	dbConfigs := repo.FindByHostName(hostname)

	for _, net := range *nets {
		indx := findHubConfigIndexByNet(dbConfigs, &net)
		if indx >= 0 {
			cfg := (*dbConfigs)[indx]
			if cfg.Ip != net.IP.String() {
				cfg.Ip = net.IP.String()
				repo.Update(cfg)
			}
		} else {
			cfg := model.HubConfig{
				HostName:   hostname,
				Interface:  net.Interface.Name,
				Ip:         net.IP.String(),
				IsMQBroker: false,
			}
			repo.InsertOne(cfg)
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
