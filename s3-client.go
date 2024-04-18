package main

import (
	"log"
	"io"
	"bytes"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type S3Client struct {
	Client *s3.S3
}

func NewS3Client() *S3Client {
	awsConfig := &aws.Config{
		Region: aws.String(Cfg.Region),
		Endpoint: aws.String(Cfg.Endpoint),
		Credentials: credentials.NewStaticCredentials(Cfg.AccessKey, Cfg.SecretKey, ""),
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
		Bucket: aws.String(Cfg.BucketName),
		Prefix: aws.String(key),
		Delimiter: aws.String("/"),
	}
	return s.Client.ListObjectsV2(input)
}

func (s *S3Client) GetObject(key string) (*s3.GetObjectOutput, error) {
	input := &s3.GetObjectInput{
		Bucket: aws.String(Cfg.BucketName),
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
		Bucket: aws.String(Cfg.BucketName),
		Key:    aws.String(key),
		Body:   bytes.NewReader(data),
	}
	return s.Client.PutObject(input)
}

func (s *S3Client) DeleteObject(key string) (*s3.DeleteObjectOutput, error) {
	input := &s3.DeleteObjectInput{
		Bucket: aws.String(Cfg.BucketName),
		Key: aws.String(key),
	}
	return s.Client.DeleteObject(input)
}

func (s *S3Client) CopyObject(src, dest string) (*s3.CopyObjectOutput, error) {
	input := &s3.CopyObjectInput{
		Bucket: aws.String(Cfg.BucketName),
		CopySource: aws.String(Cfg.BucketName + "/" + src),
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
