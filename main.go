package main

import (
	// "net/http"
	"fmt"
	"log"
)

var Cfg *Config

func main() {
	var err error
	Cfg, err = LoadConfig()
	if err != nil {
		log.Fatalf("unable to load config: %v", err)
	}
	fmt.Println(Cfg.AccessKey)
	fmt.Println(Cfg.SecretKey)
	fmt.Println(Cfg.BucketName)
	fmt.Println(Cfg.Region)
	fmt.Println(Cfg.Endpoint)
	fmt.Println(Cfg.Port)
	fmt.Println(Cfg.BaseURL)
}