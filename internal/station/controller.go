package station

import (
	"context"
	"fmt"
	"iot-hub-api/internal/network"
	"iot-hub-api/internal/restclient"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

type Controller struct {
	MongoClient   *mongo.Client
	StationClient restclient.StationClient
}

func New(mongoClient *mongo.Client, stationClient restclient.StationClient) Controller {

	return Controller{
		MongoClient:   mongoClient,
		StationClient: stationClient,
	}
}

func (c *Controller) DiscoverStations(ginCtx *gin.Context, ctx context.Context) {
	netAddresses, _ := network.GetLocalAddresses()
	for _, i := range *netAddresses {
		fmt.Println(i.Interface.Name, i.IP)
		ips := network.GetAllNetworkIps(&i)
		fmt.Println(ips)
		for _, ip := range *ips {
			beaconResponse, err := c.StationClient.GetBeacon(ip.String())
			if err != nil {
				fmt.Println(beaconResponse)
			} else {
				fmt.Println(err)
			}
		}
	}

}
