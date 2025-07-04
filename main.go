package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/vitwit/matic-telemetry/client"
	"github.com/vitwit/matic-telemetry/config"
	"github.com/vitwit/matic-telemetry/rest"
)

func main() {
	cfg, err := config.ReadFromFile()
	if err != nil {
		panic(err)
	}

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	logger := log.Output(os.Stderr).With().Str("component", "telemetry").Logger()

	// Setup cancellable context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Setup graceful shutdown
	go func() {
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
		<-stop
		logger.Info().Msg("Signal received. Shutting down.")
		cancel()
	}()

	appCtx := client.NewAppContext(ctx, logger)

	err = rest.RegisterNode(appCtx, cfg)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to register node")
	}
	logger.Info().Msg("Node registered successfully")

	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	logger.Info().Msg("Telemetry client started. Sending stats every 2 seconds.")

	var retryCount int
	const maxRetries = 5
	for {
		select {
		case <-ctx.Done():
			logger.Info().Msg("Context canceled. Exiting loop.")
			return
		case <-ticker.C:
			go func() {
				err := rest.SubmitStats(appCtx, cfg)
				if err != nil {
					if err.Error() == "Node not registered" {
						retryCount++
						logger.Warn().Int("attempt", retryCount).Msg("Node not registered. Retrying registration...")

						regErr := rest.RegisterNode(appCtx, cfg)
						if regErr != nil {
							logger.Error().Err(regErr).Msg("Retry registration failed")
						} else {
							logger.Info().Msg("Node re-registered successfully")
							retryCount = 0
						}

						if retryCount >= maxRetries {
							logger.Fatal().Msg("Maximum registration retries exceeded. Shutting down.")
						}
					} else {
						logger.Error().Err(err).Msg("Failed to submit stats")
					}
				} else {
					logger.Info().Msg("Stats submitted successfully")
				}
			}()
		}
	}
}
