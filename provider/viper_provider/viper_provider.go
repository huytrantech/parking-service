package viper_provider

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
)

type IConfigProvider interface {
	GetConfigEnv() EnvConfig
}

type ConfigProvider struct {
	EnvConfig
}

type EnvConfig struct {
	ElasticSearchUrl   string `mapstructure:"ELASTIC_SEARCH_URL"`
	ElasticSearchToken string `mapstructure:"ELASTIC_SEARCH_TOKEN"`
	GoogleMapDomain    string `mapstructure:"GOOGLE_MAP_DOMAIN"`
	GoogleMapKey       string `mapstructure:"GOOGLE_MAP_KEY"`
	RedisAddress       string `mapstructure:"REDIS_ADDRESS"`
	PGDatabase         string `mapstructure:"PG_DATABASE"`
	RabbitMQUrl        string `mapstructure:"RABBIT_MQ_URL"`
	GoongUrl           string `mapstructure:"GOONG_URL"`
	GoongKey           string `mapstructure:"GOONG_KEY"`
}

func NewConfigProvider() (configProvider IConfigProvider) {

	envCf := EnvConfig{}
	prefixEnv := "local"
	if os.Getenv("ENV") == "prod" {
		prefixEnv = os.Getenv("ENV")
	}
	viper.AddConfigPath(".")
	viper.SetConfigName(fmt.Sprintf("environment.%s.env", prefixEnv))
	viper.SetConfigType("env")
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	err = viper.Unmarshal(&envCf)

	return &ConfigProvider{EnvConfig: envCf}
}

func (cf *ConfigProvider) GetConfigEnv() EnvConfig {
	return cf.EnvConfig
}
