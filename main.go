package main

import (
	"log"

	"github.com/PrathyushaLakkireddy/heimdall-node-stats/config"
	"github.com/PrathyushaLakkireddy/heimdall-node-stats/stats"
)

func main() {
	cfg, err := config.ReadFromFile()
	if err != nil {
		log.Fatal(err)
	}

	// ws://localhost:3000
	err = stats.Dailer(cfg)
	if err != nil {
		log.Printf("Error while establishing a socket connection : %v", err)
	}
}
