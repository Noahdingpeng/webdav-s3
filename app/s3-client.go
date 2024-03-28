package main

import (
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/spf13/viper"
)

type S3Client struct {
	Client *s3.S3
}

type AWSConfig struct {
	AccessKey string `mapstructure:"access_key"`
	SecretKey string `mapstructure:"secret_key"`
	BucketName string `mapstructure:"bucket_name"`
	Region string `mapstructure:"region"`
	Endpoint string `mapstructure:"endpoint"`
}

func newS3Config() string {
	path := "../conf/"
	_, err := os.Stat(path+"config.yaml")
	if err != nil {
		fi, err := os.ReadFile("config_sample.yaml")
		if err != nil {
			log.Fatalf("unable to read sample config: %v", err)
		}
		err = os.WriteFile(path+"config.yaml", fi, 0644)
		if err != nil {
			log.Fatalf("unable to write config: %v", err)
		}
	}
	return path
}

func loadConfig(path string) (config AWSConfig, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil{
		log.Fatalf("unable to read config: %v, %v", path, err)
	}
	err = viper.Unmarshal(&config)
	if err != nil {
		log.Fatalf("unable to unmarshal config: %v", err)
	}
	return config, err
}

func NewS3Client() *S3Client {
	config, err := loadConfig(newS3Config())
	if err != nil {
		log.Fatalf("unable to load config: %v", err)
	}
	awsConfig := &aws.Config{
		Region: aws.String(config.Region),
		Endpoint: aws.String(config.Endpoint),
		Credentials: credentials.NewStaticCredentials(config.AccessKey, config.SecretKey, ""),
	}
	awsConfig.S3ForcePathStyle = aws.Bool(true)

	sess, err := session.NewSession(awsConfig)
	if err != nil {
		log.Fatalf("unable to create session: %v", err)
	}
	s3Client := s3.New(sess)
	return &S3Client{
		Client: s3Client,
	}
}
