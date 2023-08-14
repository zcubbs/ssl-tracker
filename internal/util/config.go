package util

import (
	"github.com/charmbracelet/log"
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
	log.Info("loaded configuration")
	return config
}

func loadConfig() {

	err := godotenv.Load(".env")

	if err != nil {
		if viper.GetString("debug") == "true" {
			log.Warn("no .env file found")
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
	if err != nil && viper.GetString("debug") == "true" {
		log.Warn("unable to load config file",
			"path", ViperConfigName+"."+ViperConfigType,
		)
	}

	viper.AutomaticEnv() // read in environment variables that match
	viper.SetEnvPrefix(ViperConfigEnvPrefix)

	for _, key := range allowedEnvVarKeys {
		err := viper.BindEnv(key)
		if err != nil {
			log.Printf("error: %s", err)
		}
	}
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	err = viper.Unmarshal(&config)
	if err != nil {
		log.Printf("warn: could not decode config into struct: %v", err)
	}
}
