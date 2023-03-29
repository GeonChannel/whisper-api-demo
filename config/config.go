package config

import (
	"fmt"
	"github.com/spf13/viper"
)

func SetViper() {
	// Set config type
	viper.SetConfigType("json")

	// Set config file name and path
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")

	// Read config file
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("failed to read config file: %s", err))
	}
}
