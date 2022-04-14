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
		//fmt.Println(ips)
		fmt.Println("Looking for alive stations")
		for _, ip := range *ips {
			beaconResponse, _ := c.StationClient.GetBeacon(ip.String())
			/*if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println(*beaconResponse)
			}*/
			if beaconResponse != nil {
				fmt.Println(beaconResponse)
			}
		}
		fmt.Println("END Looking for alive stations")
	}

}
