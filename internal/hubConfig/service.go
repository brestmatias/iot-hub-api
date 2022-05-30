package hub_config

import "iot-hub-api/internal/repository"

type HubConfigService struct {
	HubConfigRepository *repository.HubConfigRepository
}

func NewHubConfigService(hubConfigRepository *repository.HubConfigRepository) *HubConfigService {
	return &HubConfigService{
		HubConfigRepository: hubConfigRepository,
	}
}

func (s *HubConfigService) GetBrokerAddress() string {
	configs := (*s.HubConfigRepository).FindByField("is_broker", true)
	if configs != nil && len(*configs) > 0 {
		return (*configs)[0].Ip
	}
	return ""
}
