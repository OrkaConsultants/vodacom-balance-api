package main

import (
	log "github.com/sirupsen/logrus"

	"github.com/OrkaConsultants/vodacom-balance-api/internal/core"
)

func main() {
	log.SetLevel(log.InfoLevel)

	// Load all the required config files
	core.LoadConfig()

	// Start up the API server
	core.SetupServer()

	// Register with Eureka server
	go core.SetupEureka()

	select {} // block forever
}
