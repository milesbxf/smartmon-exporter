package collector

import (
	"github.com/milesbxf/smartmon-exporter/pkg/smartctl"
	"github.com/prometheus/client_golang/prometheus"
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
}

func (c *collector) Describe(descs chan<- *prometheus.Desc) {
	for _, m := range c.metrics {
		m.Describe(descs)
	}
}

func (c *collector) Collect(metrics chan<- prometheus.Metric) {
	for _, m := range c.metrics {
		m.Collect(metrics)
	}
}

func (c *collector) Run() error {
	c.poll()
	t := time.NewTicker(c.pollInterval)
	for range t.C {
		c.poll()
	}
	return nil
}

func (c *collector) poll() error {
	for i, d := range c.devices {
		info, err := c.smart.InfoAll(d)
		if err != nil {
			// TODO: better error handling
			return err
		}
		for _, m := range c.metrics[i].metrics {
			err := m.UpdateFromInfo(*info)
			if err != nil {
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
