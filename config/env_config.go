package config

import (
	"github.com/spf13/viper"
	"log"
	"os"
)

type EnvConfig struct {
	Test        string
	Deploy      string
	ServicesUrl string `json:"services_url"`
}

var Configuration EnvConfig

func init() {
	env := os.Getenv("GO_ENV")
	if env == "" {
		env = "dev"
	}

	viper.SetConfigName("config/config." + env)
	viper.SetConfigType("json")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Fatal error config file: %s \n", err)
	}

	if err := viper.Unmarshal(&Configuration); err != nil {
		log.Fatalf("Unable to decode into struct: %s \n", err)
	}
}

func GetConfig() EnvConfig {
	return Configuration
}
