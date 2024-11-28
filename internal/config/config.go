package config

import "github.com/spf13/viper"

type Config struct {
	Service string
	Env     string

	Debug bool

	DB struct {
		DSN string `yaml:"dsn"`
	} `yaml:"db"`
}

func Load(service, env string) (*Config, error) {
	viper.SetConfigName(env)
	viper.SetConfigType("yaml")
	viper.AddConfigPath("configs/")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	cfg.Service = service
	cfg.Env = env

	return &cfg, nil
}
