package collector

import (
	"github.com/milesbxf/smartmon-exporter/pkg/smartctl"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog/log"
	"sync"
	"time"
)

type Collector interface {
	Run() error
}

type collector struct {
	smart        smartctl.SmartCtl
	devices      []string
	metrics      []*Metrics
	pollInterval time.Duration
	mu           sync.RWMutex
}

func (c *collector) Describe(descs chan<- *prometheus.Desc) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	for _, m := range c.metrics {
		m.Describe(descs)
	}
}

func (c *collector) Collect(metrics chan<- prometheus.Metric) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	for _, m := range c.metrics {
		m.Collect(metrics)
	}
}

func (c *collector) Run() error {
	if err := c.poll(); err != nil {
		log.Error().Err(err).Msg("failed to do initial poll")
	}

	t := time.NewTicker(c.pollInterval)
	for range t.C {
		if err := c.poll(); err != nil {
			log.Error().Err(err).Msg("failed to poll")
		}
	}
	return nil
}

func (c *collector) poll() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	for i, d := range c.devices {
		info, err := c.smart.InfoAll(d)
		if err != nil {
			return err
		}
		log.Info().Str("device", d).Msg("got info")

		for _, m := range c.metrics[i].metrics {
			if err := m.UpdateFromInfo(*info); err != nil {
				return err
			}
		}
	}
	return nil
}

func New(smart smartctl.SmartCtl, pollInterval time.Duration) (*collector, error) {
	scan, err := smart.ScanOpen()
	if err != nil {
		return nil, err
	}

	devices := []string{}
	metrics := []*Metrics{}

	for _, d := range scan.Devices {
		log.Info().Str("device", d.Name).Msg("found device")
		devices = append(devices, d.Name)
		metrics = append(metrics, NewMetrics())
	}

	return &collector{
		smart:        smart,
		devices:      devices,
		metrics:      metrics,
		pollInterval: pollInterval,
	}, nil
}
