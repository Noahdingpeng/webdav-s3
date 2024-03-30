package main

import (
	"log"
	"os"
	"io"
	"bytes"

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

var Config AWSConfig

func newS3Config() string {
	path := "conf/"
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
	err = viper.Unmarshal(&Config)
	if err != nil {
		log.Fatalf("unable to unmarshal config: %v", err)
	}
	return Config, err
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

func (s *S3Client) ListObjects(key string) (*s3.ListObjectsV2Output, error) {
	input := &s3.ListObjectsV2Input{
		Bucket: aws.String(Config.BucketName),
		Prefix: aws.String(key),
		Delimiter: aws.String("/"),
	}
	return s.Client.ListObjectsV2(input)
}

func (s *S3Client) GetObject(key string) (*s3.GetObjectOutput, error) {
	input := &s3.GetObjectInput{
		Bucket: aws.String(Config.BucketName),
		Key: aws.String(key),
	}
	return s.Client.GetObject(input)
}

func (s *S3Client) PutObject(key string, body io.Reader) (*s3.PutObjectOutput, error) {
	data, err := io.ReadAll(body)
	if err != nil {
		return nil, err
	}
	input := &s3.PutObjectInput{
		Bucket: aws.String(Config.BucketName),
		Key:    aws.String(key),
		Body:   bytes.NewReader(data),
	}
	return s.Client.PutObject(input)
}

func (s *S3Client) DeleteObject(key string) (*s3.DeleteObjectOutput, error) {
	input := &s3.DeleteObjectInput{
		Bucket: aws.String(Config.BucketName),
		Key: aws.String(key),
	}
	return s.Client.DeleteObject(input)
}

func (s *S3Client) CopyObject(src, dest string) (*s3.CopyObjectOutput, error) {
	input := &s3.CopyObjectInput{
		Bucket: aws.String(Config.BucketName),
		CopySource: aws.String(Config.BucketName + "/" + src),
		Key: aws.String(dest),
	}
	return s.Client.CopyObject(input)
}

func (s *S3Client) MoveObject(src, dest string) (*s3.CopyObjectOutput, error) {
	_, err := s.CopyObject(src, dest)
	if err != nil {
		return nil, err
	}
	_, err = s.DeleteObject(src)
	if err != nil {
		return nil, err
	}
	return nil, nil
}
