package config

import (
	"encoding/json"
	"log"
	"os"

	"github.com/spf13/viper"
)

func LoadConfig(path *string) (*Config, error) {

	fileEnv := "../../externalConnectionStrings.env.development"

	if _, err := os.Stat(fileEnv); err != nil {
		if os.IsNotExist(err) {
			fileEnv = "../../externalConnectionStrings.env"
		} else {
			log.Fatalf("Erro ao verificar arquivo %s: %v", fileEnv, err)
		}
	}

	viper.SetConfigFile(fileEnv)
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Erro ao ler .env: %v", err)
	}

	file, err := os.Open(*path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var cfg Config
	if err := json.NewDecoder(file).Decode(&cfg.File); err != nil {
		return nil, err
	}

	if cfg.Server == nil {
		cfg.Server = &ServerConfig{}
	}

	if c := viper.GetString("DB_CONNECTION"); c != "" {
		cfg.Server.ConnectionString = c
	}
	if t := viper.GetString("DB_TYPE"); t != "" {
		cfg.Server.DbType = t
	}

	return &cfg, nil
}

func OverrideConfig(cfg *Config, conectionString, dbType, file, outDir, procsFlag *string, workersNum *int) {
	if *conectionString != "" {
		cfg.Server.ConnectionString = *conectionString
	}
	if *dbType != "" {
		cfg.Server.DbType = *dbType
	}
	if *file != "" {
		cfg.File.File = *file
	}
	if *outDir != "" {
		cfg.File.Out = *outDir
	}
	if *procsFlag != "" {
		cfg.File.Procs = *procsFlag
	}
	if *workersNum > 0 {
		cfg.File.NumWorkers = *workersNum
	}
	if cfg.File.NumWorkers == 0 {
		cfg.File.NumWorkers = 5
	}
}
