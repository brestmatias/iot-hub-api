package app

import (
	"context"
	"fmt"
	"iot-hub-api/internal/config"
	"iot-hub-api/internal/restclient"
	"iot-hub-api/internal/station"
	"log"
	"time"

	"github.com/brestmatias/golang-restclient/rest"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// https://manytools.org/hacker-tools/ascii-banner/
const logo = `
██╗ ██████╗ ████████╗    ██╗  ██╗██╗   ██╗██████╗      █████╗ ██████╗ ██╗
██║██╔═══██╗╚══██╔══╝    ██║  ██║██║   ██║██╔══██╗    ██╔══██╗██╔══██╗██║
██║██║   ██║   ██║       ███████║██║   ██║██████╔╝    ███████║██████╔╝██║
██║██║   ██║   ██║       ██╔══██║██║   ██║██╔══██╗    ██╔══██║██╔═══╝ ██║
██║╚██████╔╝   ██║       ██║  ██║╚██████╔╝██████╔╝    ██║  ██║██║     ██║
╚═╝ ╚═════╝    ╚═╝       ╚═╝  ╚═╝ ╚═════╝ ╚═════╝     ╚═╝  ╚═╝╚═╝     ╚═╝`

type App struct {
	Configs           *config.ConfigFile
	StationController *station.Controller
}

func Start() {
	fmt.Println("⚡⚡⚡ STARTING IOT-HUB-API ⚡⚡⚡")
	fmt.Println(logo)
	ctx := context.Background()
	app := buildApp(ctx)

	if app.Configs.Server.GinMode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	buildRouter(ctx, router, app)
	fmt.Println("💪💪💪💪💪💪💪💪💪 IOT-HUB-API READY STARTING SERVER IN PORT: ", app.Configs.Server.Port)
	router.Run(":" + app.Configs.Server.Port)
}

func buildApp(ctx context.Context) *App {
	configs := config.GetConfigs()

	mongoClient := buildMongoClient(configs, ctx)
	stationClient := buildRestClients(configs)

	stationController := station.New(mongoClient, stationClient)

	return &App{
		Configs:           configs,
		StationController: &stationController,
	}
}

func buildMongoClient(configs *config.ConfigFile, ctx context.Context) *mongo.Client {
	clientOptions := options.Client().ApplyURI(configs.Database.Uri)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	return client
}

func buildRestClients(configs *config.ConfigFile) restclient.StationClient {
	customPool := &rest.CustomPool{
		MaxIdleConnsPerHost: 100,
	}

	rb := rest.RequestBuilder{
		BaseURL:        configs.StationRestClient.BaseURL,
		ConnectTimeout: time.Duration(configs.StationRestClient.ConnectTimeout) * time.Millisecond,
		Timeout:        time.Duration(configs.StationRestClient.Timeout) * time.Millisecond,
		ContentType:    rest.JSON,
		DisableCache:   configs.StationRestClient.DisableCache,
		DisableTimeout: configs.StationRestClient.DisableTimeout,
		CustomPool:     customPool,
	}

	return restclient.NewStationClient(&rb)
}
