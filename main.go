package main

import (
	"github.com/plumbie/plumbie/config"
	"github.com/plumbie/plumbie/models"
	"github.com/plumbie/plumbie/webserver"

	log "github.com/sirupsen/logrus"
)

func main() {
	if err := config.LoadConfig(); err != nil {
		log.Panic(err)
	}

	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
	if config.Debug {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	if err := models.Initialize(); err != nil {
		log.Panic(err)
	}
	webserver.Start()
}
