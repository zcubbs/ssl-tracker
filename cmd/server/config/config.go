package config

import (
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	lwl "github.com/zcubbs/logwrapper/logger"
	"github.com/zcubbs/ssl-tracker/cmd/server/logger"

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
	log = logger.L()
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
		logger.SetFormat(lwl.JSONFormat)
		logger.SetLevel(lwl.InfoLevel)
		log.Info("production=true")
	} else {
		logger.SetFormat(lwl.TextFormat)
		logger.SetLevel(lwl.DebugLevel)
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

	for _, key := range allowedEnvVarKeys {
		err := viper.BindEnv(key)
		if err != nil {
			log.Error("unable to bind environment variable", "key", key, "error", err)
		}
	}
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	err = viper.Unmarshal(&cfg)
	if err != nil {
		log.Warn("unable to unmarshal config", "error", err)
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
