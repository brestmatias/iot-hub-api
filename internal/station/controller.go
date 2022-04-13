package station

import (
	"context"
	"fmt"
	"log"
	"net"

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
		for _, a := range addrs {
			log.Printf("%v %v\n", i.Name, a)
			a.Network()
		}
	}
}
