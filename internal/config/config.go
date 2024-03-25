package config

import (
	"github.com/spf13/viper"
)

type Global struct {
	Gin *Gin
}

type Gin struct {
	Port  string
	Mysql *Mysql
	Redis *Redis
	Token *Token
}
type Mysql struct {
	DSN string
}

type Redis struct {
	Addr     string
	Password string
	DB       int
}
type Token struct {
	PrivateKeyPath string
	ValidSeconds   int
}

func NewConfig(configPath, configName string) (*Global, error) {
	var config Global
	viper.SetConfigName(configName)
	viper.AddConfigPath(configPath)
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}
	viper.Unmarshal(&config)
	return &config, nil
}
