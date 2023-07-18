package util

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"strings"
	"sync"
)

var (
	config     Config
	onceConfig sync.Once
)

func Bootstrap() Config {
	onceConfig.Do(loadConfig)
	return config
}

func loadConfig() {

	err := godotenv.Load(".env")

	if err != nil {
		if viper.GetString("debug") == "true" {
			fmt.Println("loading .env file")
		}
	}

	for k, v := range defaults {
		viper.SetDefault(k, v)
	}

	for _, p := range viperConfigPaths {
		viper.AddConfigPath(p)
	}

	viper.SetConfigType(ViperConfigType)
	viper.SetConfigName(ViperConfigName)

	err = viper.ReadInConfig()
	if err != nil {
		fmt.Printf("warn: %s\n", err)
	}

	viper.AutomaticEnv() // read in environment variables that match
	viper.SetEnvPrefix(ViperConfigEnvPrefix)

	for _, key := range allowedEnvVarKeys {
		err := viper.BindEnv(key)
		if err != nil {
			fmt.Printf("error: %s", err)
		}
	}
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	err = viper.Unmarshal(&config)
	if err != nil {
		fmt.Printf("warn: could not decode config into struct: %v", err)
	}
}
