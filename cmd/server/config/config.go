package config

import (
	"encoding/json"
	"fmt"
	"github.com/charmbracelet/log"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"os"
	"strings"
	"sync"
)

var (
	cfg        Config
	onceEnv    sync.Once
	onceLogger sync.Once
	onceConfig sync.Once
)

var (
	Version string
	Commit  string
	Date    string
)

func Bootstrap() Config {
	onceEnv.Do(LoadEnv)
	onceLogger.Do(setupLogger)
	onceConfig.Do(loadConfig)
	log.Info("loaded configuration")
	return cfg
}

func LoadEnv() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Debug("no .env file found")
	}
}

func setupLogger() {
	// read environment variable for log level
	env := os.Getenv("ENVIRONMENT")
	if env == "production" || env == "prod" {
		log.SetLevel(log.InfoLevel)
		log.SetFormatter(log.JSONFormatter)
	} else {
		log.SetLevel(log.DebugLevel)
	}
}

func loadConfig() {
	for k, v := range defaults {
		viper.SetDefault(k, v)
	}

	for _, p := range viperConfigPaths {
		viper.AddConfigPath(p)
	}

	viper.SetConfigType(ViperConfigType)
	viper.SetConfigName(ViperConfigName)

	err := viper.ReadInConfig()
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

	err = viper.Unmarshal(&cfg)
	if err != nil {
		log.Printf("warn: could not decode config into struct: %v", err)
	}

	if viper.GetString("debug") == "true" {
		debugConfig()
	}
}

func debugConfig() {
	jsonConfig, err := json.MarshalIndent(&cfg, "", "  ")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%v\n", string(jsonConfig))
}
