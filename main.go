//go:generate swagger generate spec /o static/spec.json

// Package api Crawler API.
//
// Web Crawler API
//     Schemes: http
//     Host: localhost
//     Version: 0.1
//     Contact: foreverqihe@gmail.com
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
// swagger:meta
package main

import (
	"os"

	"github.com/foreverqihe/crawler/api"
	"github.com/foreverqihe/crawler/core"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
}

func main() {
	log.Infoln("Crawler started")
	crawler := api.CrawlerAPI{Handler: core.NewHander()}
	if err := crawler.Run(); err != nil {
		log.Fatalln("failed to start api", err)
	}
}
