package app

import (
	"log"

	"iot-hub-api/internal/network"
	"iot-hub-api/internal/repository"
	"iot-hub-api/model"
)

func MapCurrentHostInterfaces(repo repository.HubConfigRepository) {
	hostname := network.GetHostName()
	log.Println("Local Hostname:", hostname)
	nets, _ := network.GetLocalAddresses()
	log.Println("Local Networks:", nets)

	dbConfigs := repo.FindAll()
	log.Println("DBInterfaces configs:", dbConfigs)

	log.Println("Merging Local Interfaces With DB")
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
				Mac: net.Interface.HardwareAddr.String(),
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
		if dbCfg.Mac == net.Interface.HardwareAddr.String() {
			return i
		}
	}
	return -1
}
