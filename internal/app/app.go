package app

import (
	"context"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// https://manytools.org/hacker-tools/ascii-banner/
const logo = `
██▓ ▒█████  ▄▄▄█████▓   ▓█████▄  ██▓  ██████  ██▓███   ▄▄▄     ▄▄▄█████▓ ▄████▄   ██░ ██ ▓█████  ██▀███  
▓██▒▒██▒  ██▒▓  ██▒ ▓▒   ▒██▀ ██▌▓██▒▒██    ▒ ▓██░  ██▒▒████▄   ▓  ██▒ ▓▒▒██▀ ▀█  ▓██░ ██▒▓█   ▀ ▓██ ▒ ██▒
▒██▒▒██░  ██▒▒ ▓██░ ▒░   ░██   █▌▒██▒░ ▓██▄   ▓██░ ██▓▒▒██  ▀█▄ ▒ ▓██░ ▒░▒▓█    ▄ ▒██▀▀██░▒███   ▓██ ░▄█ ▒
░██░▒██   ██░░ ▓██▓ ░    ░▓█▄   ▌░██░  ▒   ██▒▒██▄█▓▒ ▒░██▄▄▄▄██░ ▓██▓ ░ ▒▓▓▄ ▄██▒░▓█ ░██ ▒▓█  ▄ ▒██▀▀█▄  
░██░░ ████▓▒░  ▒██▒ ░    ░▒████▓ ░██░▒██████▒▒▒██▒ ░  ░ ▓█   ▓██▒ ▒██▒ ░ ▒ ▓███▀ ░░▓█▒░██▓░▒████▒░██▓ ▒██▒
░▓  ░ ▒░▒░▒░   ▒ ░░       ▒▒▓  ▒ ░▓  ▒ ▒▓▒ ▒ ░▒▓▒░ ░  ░ ▒▒   ▓▒█░ ▒ ░░   ░ ░▒ ▒  ░ ▒ ░░▒░▒░░ ▒░ ░░ ▒▓ ░▒▓░
 ▒ ░  ░ ▒ ▒░     ░        ░ ▒  ▒  ▒ ░░ ░▒  ░ ░░▒ ░       ▒   ▒▒ ░   ░      ░  ▒    ▒ ░▒░ ░ ░ ░  ░  ░▒ ░ ▒░
 ▒ ░░ ░ ░ ▒    ░          ░ ░  ░  ▒ ░░  ░  ░  ░░         ░   ▒    ░      ░         ░  ░░ ░   ░     ░░   ░ 
 ░      ░ ░                 ░     ░        ░                 ░  ░        ░ ░       ░  ░  ░   ░  ░   ░     
                          ░                                              ░                                `

type App struct {
	Configs           *config.Configs
	StationController *station.Controller
}

func Start() {
	ctx := context.Background()
	app := buildApp(ctx)

	if app.Configs.GinMode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	buildRouter(ctx, router, app)
	router.Run(":" + app.Configs.Port)
	fmt.Println("💪💪💪💪💪💪💪💪💪 IOT-DISPATCHER READY AND UP SERVING IN PORT: ", app.Configs.Port)
}

func buildApp(ctx context.Context) *App {
	fmt.Println("⚡⚡⚡ STARTING IOT-DISPATCHER ⚡⚡⚡")
	fmt.Println(logo)
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
