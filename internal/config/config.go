package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type Configs struct {
	Verbose bool
	Port    string
	GinMode string

	DBUri string
}

type ConfigFile struct {
	Server struct {
		Port    string `yaml:"port"`
		Verbose int    `yaml:"verbose"`
		GinMode string `yaml:"ginMode"`
	} `yaml:"server"`
	Database struct {
		Uri string `yaml:"uri"`
	} `yaml:"database"`
}

func GetConfigs() *Configs {
	configs := Configs{}
	env := os.Getenv("IOTENV")
	if env == "" {
		env = "dev"
	}
	filename := "config." + env + ".yml"
	f, err := os.Open(filename)
	if err != nil {
		log.Println("Error abriendo archivo de configuracion, " + filename + " - " + err.Error())
	}
	defer f.Close()

	var cfgFile ConfigFile
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfgFile)
	if err != nil {
		log.Println("Error abriendo archivo de configuracion, " + err.Error())
	}
	log.Println("GET-CONFIG-FOR", env)

	if ginMode := os.Getenv("GIN_MODE"); ginMode == "release" {
		configs.GinMode = "release"
	} else {
		configs.GinMode = cfgFile.Server.GinMode
	}

	configs.Port = cfgFile.Server.Port
	if cfgFile.Server.Verbose == 1 {
		configs.Verbose = true
	}

	configs.DBUri = cfgFile.Database.Uri
	return &configs
}
