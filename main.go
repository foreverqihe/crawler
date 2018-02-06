package main

import (
	"github.com/foreverqihe/crawler/api"
	log "github.com/sirupsen/logrus"
	"os"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
}

func main() {
	log.Infoln("Crawler started")
	crawler := new(api.CrawlerApi)
	if err := crawler.Run(); err != nil {
		log.Fatalln("failed to start api", err)
	}
}
