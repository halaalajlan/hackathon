package main

import (
	"os"
	"os/signal"

	"github.com/halaalajlan/hackathon/api"
	log "github.com/halaalajlan/hackathon/logger"
	"github.com/halaalajlan/hackathon/models"
)

func main() {
	err := models.SetUp()
	if err != nil {
		log.Error(err)
	}
	server := api.NewServer()

	go server.Start()

	// Handle graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	log.Info("CTRL+C Received... Gracefully shutting down servers")
	server.Shutdown()

}
