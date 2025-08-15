package config

import (
	"encoding/json"
	"os"

	"github.com/spf13/viper"
)

func Load(path string) (*Config, error) {
	viper.AutomaticEnv()

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var cfg Config
	if err := json.NewDecoder(file).Decode(&cfg); err != nil {
		return nil, err
	}

	if s := viper.GetString("DB_SERVER"); s != "" {
		if cfg.Server == nil {
			cfg.Server = &ServerConfig{}
		}
		cfg.Server.Server = s
	}
	if u := viper.GetString("DB_USER"); u != "" {
		if cfg.Server == nil {
			cfg.Server = &ServerConfig{}
		}
		cfg.Server.User = u
	}
	if p := viper.GetString("DB_PASS"); p != "" {
		if cfg.Server == nil {
			cfg.Server = &ServerConfig{}
		}
		cfg.Server.Password = p
	}
	if d := viper.GetString("DB_NAME"); d != "" {
		if cfg.Server == nil {
			cfg.Server = &ServerConfig{}
		}
		cfg.Server.Database = d
	}

	return &cfg, nil
}
