package main

import (
	"log"

	"github.com/PrathyushaLakkireddy/heimdall-node-stats/config"
	"github.com/PrathyushaLakkireddy/heimdall-node-stats/stats"
)

func main() {
	// Checking for config
	cfg, err := config.ReadFromFile()
	if err != nil {
		log.Fatal(err)
	}

	// Calling dailer to establish connection with netstats and
	// report metrics to the dashboard
	err = stats.Dailer(cfg)
	if err != nil {
		log.Printf("Error while establishing a socket connection : %v", err)
	}
}
