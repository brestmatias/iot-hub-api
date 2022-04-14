package station

import (
	"context"
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

type Controller struct {
	MongoClient *mongo.Client
}

func New(mongoClient *mongo.Client) Controller {

	return Controller{
		MongoClient: mongoClient,
	}
}

func (controller *Controller) DiscoverStations(ginCtx *gin.Context, ctx context.Context) {
	localAddresses()

}

func localAddresses() {
	ifaces, err := net.Interfaces()
	if err != nil {
		log.Print(fmt.Errorf("localAddresses: %v\n", err.Error()))
		return
	}

	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			log.Print(fmt.Errorf("localAddresses: %v\n", err.Error()))
			continue
		}
		if strings.Contains(i.Flags.String(), "up") {
			for _, a := range addrs {
				ip, net, _ := net.ParseCIDR(a.String())
				fmt.Printf("%v %v %v %v %v %v\n", i.Name, a, ip.To4(), ip, net.IP.IsLoopback(), ip.IsPrivate())

			}
		}

	}

}
