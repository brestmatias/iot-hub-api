package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

// OJO!!! IDENTADO DEBE HACERSE CON ESPACIOS EN EL YML
type ConfigFile struct {
	Server struct {
		Port    string `yaml:"port"`
		GinMode string `yaml:"ginMode"`
	} `yaml:"server"`
	Database struct {
		Uri string `yaml:"uri"`
		DB  string `yaml:"db"`
	} `yaml:"database"`
	StationRestClient struct {
		BaseURL        string `yaml:"baseURL"`
		ConnectTimeout int    `yaml:"connectTimeout"`
		Timeout        int    `yaml:"timeout"`
		DisableCache   bool   `yaml:"disableCache"`
		DisableTimeout bool   `yaml:"disableTimeout"`
	} `yaml:"station-rest-client"`
	SlowStationRestClient struct {
		BaseURL        string `yaml:"baseURL"`
		ConnectTimeout int    `yaml:"connectTimeout"`
		Timeout        int    `yaml:"timeout"`
		DisableCache   bool   `yaml:"disableCache"`
		DisableTimeout bool   `yaml:"disableTimeout"`
	} `yaml:"slow-station-rest-client"`
	Mqtt struct {
		MinInterval         string `yaml:"minInterval"` // Min interval to send same previous menssage 
		StationCommandTopic string `yaml:"stationCommandTopic"`
	} `yaml:"mqtt"`
}

func GetConfigs() *ConfigFile {
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
		cfgFile.Server.GinMode = "release"
	}

	return &cfgFile
}
