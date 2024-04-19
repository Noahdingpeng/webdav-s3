package main

import (
	"log"
	"strings"
)

func Logoutput(message string, level string) {
	switch level {
		case "debug":
			if strings.ToLower(Cfg.Loglevel) == "debug" {
				log.Println("DEBUG: " + message)
			}else{
				return
			}
		case "info":
			if strings.ToLower(Cfg.Loglevel) == "debug" || strings.ToLower(Cfg.Loglevel) == "info" {
				log.Println("INFO: " + message)
			}else{
				return
			}
		case "error":
			log.Fatalln("ERROR: " + message)
		default:
			return
	}
}
