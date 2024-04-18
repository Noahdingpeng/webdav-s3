package main

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	AccessKey string `mapstructure:"access_key"`
	SecretKey string `mapstructure:"secret_key"`
	BucketName string `mapstructure:"bucket_name"`
	Region string `mapstructure:"region"`
	Endpoint string `mapstructure:"endpoint"`
	
	Port string `mapstructure:"port"`
	BaseURL string `mapstructure:"base_url"`
}

func LoadConfig() (*Config, error) {
	viper.SetConfigName("config_sample")
	viper.AddConfigPath(".")
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("unable to read sample config: %v", err)
		return nil, err
	}

	viper.AutomaticEnv()

	viper.SetConfigName("config")
	viper.AddConfigPath("conf")
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		log.Println("unable to read config, using Environment Variables")
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("unable to unmarshal config: %v", err)
		return nil, err
	}
	return &config, nil
}
