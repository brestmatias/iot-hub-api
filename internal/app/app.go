package app

import (
	"context"
	"fmt"
	"iot-hub-api/internal/config"
	"iot-hub-api/internal/station"
	"log"

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
	Configs           *config.Configs
	StationController *station.Controller
}

func Start() {
	fmt.Println("⚡⚡⚡ STARTING IOT-HUB-API ⚡⚡⚡")
	fmt.Println(logo)
	ctx := context.Background()
	app := buildApp(ctx)

	if app.Configs.GinMode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	buildRouter(ctx, router, app)
	fmt.Println("💪💪💪💪💪💪💪💪💪 IOT-HUB-API READY STARTING SERVER IN PORT: ", app.Configs.Port)
	router.Run(":" + app.Configs.Port)
}

func buildApp(ctx context.Context) *App {
	configs := config.GetConfigs()

	client := buildMongoClient(configs, ctx)

	stationController := station.New(client)

	return &App{
		Configs:           configs,
		StationController: &stationController,
	}
}

func buildMongoClient(configs *config.Configs, ctx context.Context) *mongo.Client {
	clientOptions := options.Client().ApplyURI(configs.DBUri)
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
