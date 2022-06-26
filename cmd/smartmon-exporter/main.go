package main

import (
	"flag"
	"github.com/milesbxf/smartmon-exporter/pkg/collector"
	"github.com/milesbxf/smartmon-exporter/pkg/smartctl"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"net/http"
	"time"
)

func main() {
	addr := flag.String("listen-address", ":9101", "The address to listen on for HTTP requests.")
	pollIntervalStr := flag.String("poll-interval", "1m", "The interval between polling for device information.")
	flag.Parse()

	pollInterval, err := time.ParseDuration(*pollIntervalStr)
	if err != nil {
		log.Fatal().Err(err).Msgf("Could not parse poll interval %s", *pollIntervalStr)
	}

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	s := smartctl.New()

	c, err := collector.New(s, pollInterval)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create collector")
	}

	go func() {
		err := c.Run()
		if err != nil {
			log.Fatal().Err(err).Msg("failed to run collector")
		}
	}()

	if err := prometheus.Register(c); err != nil {
		log.Fatal().Err(err).Msg("failed to register collector")
	}
	http.Handle("/metrics", promhttp.Handler())

	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal().Err(err).Msg("failed to start http server")
	}
}
