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
	flag.Parse()
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	s := smartctl.New()

	c, err := collector.New(s, time.Minute)
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
