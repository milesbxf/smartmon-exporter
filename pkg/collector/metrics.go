package collector

import (
	"github.com/milesbxf/smartmon-exporter/pkg/smartctl"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog/log"
)

type PerDeviceInfoMetric interface {
	Desc() *prometheus.Desc
	Update(chan<- prometheus.Metric) error
	UpdateFromInfo(info smartctl.InfoAllOutput) error
}

type Metrics struct {
	metrics []PerDeviceInfoMetric
}

func (m *Metrics) Describe(descs chan<- *prometheus.Desc) {
	for _, m := range m.metrics {
		descs <- m.Desc()
	}
}

func (m Metrics) Collect(metrics chan<- prometheus.Metric) {
	log.Info().Msg("collecting all metrics")
	for _, m := range m.metrics {
		_ = m.Update(metrics)
	}
	log.Info().
		Int("num_metrics", len(m.metrics)).
		Msg("collected all metrics")
}
func (m *Metrics) UpdateFromInfo(info smartctl.InfoAllOutput) error {
	for _, m := range m.metrics {
		if err := m.UpdateFromInfo(info); err != nil {
			return err
		}
	}
	return nil
}

func NewMetrics() *Metrics {
	return &Metrics{
		metrics: metrics,
	}
}

type infoMetric struct {
	PromDesc    *prometheus.Desc
	UpdateFunc  func(chan<- prometheus.Metric, smartctl.InfoAllOutput, *prometheus.Desc) error
	lastInfo    smartctl.InfoAllOutput
	lastInfoSet bool
}

func (m *infoMetric) Desc() *prometheus.Desc {
	return m.PromDesc
}

func (m *infoMetric) Update(metrics chan<- prometheus.Metric) error {
	if m.lastInfoSet {
		return m.UpdateFunc(metrics, m.lastInfo, m.PromDesc)
	}
	return nil
}

func (m *infoMetric) UpdateFromInfo(info smartctl.InfoAllOutput) error {
	m.lastInfo = info
	m.lastInfoSet = true
	return nil
}
