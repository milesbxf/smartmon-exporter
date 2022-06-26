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
		metrics: []PerDeviceInfoMetric{
			&infoMetric{
				PromDesc: prometheus.NewDesc(
					"smart_device_info",
					"Information about the device",
					[]string{"device", "model_family", "model_name", "serial_number", "firmware_version"},
					nil,
				),
				UpdateFunc: func(metrics chan<- prometheus.Metric, output smartctl.InfoAllOutput, desc *prometheus.Desc) error {
					metrics <- prometheus.MustNewConstMetric(
						desc,
						prometheus.GaugeValue,
						1,
						output.Device.Name,
						output.ModelFamily,
						output.ModelName,
						output.SerialNumber,
						output.FirmwareVersion,
					)
					return nil
				},
			},
			&infoMetric{
				PromDesc: prometheus.NewDesc(
					"smart_device_user_capacity_blocks",
					"User capacity of the device",
					[]string{"device"},
					nil,
				),
				UpdateFunc: func(metrics chan<- prometheus.Metric, output smartctl.InfoAllOutput, desc *prometheus.Desc) error {
					metrics <- prometheus.MustNewConstMetric(
						desc,
						prometheus.GaugeValue,
						float64(output.UserCapacity.Blocks),
						output.Device.Name,
					)
					return nil
				},
			},
			&infoMetric{
				PromDesc: prometheus.NewDesc(
					"smart_device_user_capacity_bytes",
					"User capacity of the device",
					[]string{"device"},
					nil,
				),
				UpdateFunc: func(metrics chan<- prometheus.Metric, output smartctl.InfoAllOutput, desc *prometheus.Desc) error {
					metrics <- prometheus.MustNewConstMetric(
						desc,
						prometheus.GaugeValue,
						float64(output.UserCapacity.Bytes),
						output.Device.Name,
					)
					return nil
				},
			},
			&infoMetric{
				PromDesc: prometheus.NewDesc(
					"smart_device_logical_block_size",
					"Logical block size of the device",
					[]string{"device"},
					nil,
				),
				UpdateFunc: func(metrics chan<- prometheus.Metric, output smartctl.InfoAllOutput, desc *prometheus.Desc) error {
					metrics <- prometheus.MustNewConstMetric(
						desc,
						prometheus.GaugeValue,
						float64(output.LogicalBlockSize),
						output.Device.Name,
					)
					return nil
				},
			},
			&infoMetric{
				PromDesc: prometheus.NewDesc(
					"smart_device_physical_block_size",
					"Physical block size of the device",
					[]string{"device"},
					nil,
				),
				UpdateFunc: func(metrics chan<- prometheus.Metric, output smartctl.InfoAllOutput, desc *prometheus.Desc) error {
					metrics <- prometheus.MustNewConstMetric(
						desc,
						prometheus.GaugeValue,
						float64(output.PhysicalBlockSize),
						output.Device.Name,
					)
					return nil
				},
			},
			&infoMetric{
				PromDesc: prometheus.NewDesc(
					"smart_device_interface_max_speed_bits_per_second",
					"Interface speed of the device",
					[]string{"device"},
					nil,
				),
				UpdateFunc: func(metrics chan<- prometheus.Metric, output smartctl.InfoAllOutput, desc *prometheus.Desc) error {
					metrics <- prometheus.MustNewConstMetric(
						desc,
						prometheus.GaugeValue,
						float64(output.InterfaceSpeed.Max.UnitsPerSecond*output.InterfaceSpeed.Max.BitsPerUnit),
						output.Device.Name,
					)
					return nil
				},
			},
			&infoMetric{
				PromDesc: prometheus.NewDesc(
					"smart_device_smart_status_passed",
					"Whether the SMART status is a pass",
					[]string{"device"},
					nil,
				),
				UpdateFunc: func(metrics chan<- prometheus.Metric, output smartctl.InfoAllOutput, desc *prometheus.Desc) error {
					v := 0.
					if output.SmartStatus.Passed {
						v = 1.
					}
					metrics <- prometheus.MustNewConstMetric(
						desc,
						prometheus.GaugeValue,
						v,
						output.Device.Name,
					)
					return nil
				},
			},
			&infoMetric{
				PromDesc: prometheus.NewDesc(
					"smart_device_power_on_time_seconds",
					"Power on time of the device",
					[]string{"device"},
					nil,
				),
				UpdateFunc: func(metrics chan<- prometheus.Metric, output smartctl.InfoAllOutput, desc *prometheus.Desc) error {
					metrics <- prometheus.MustNewConstMetric(
						desc,
						prometheus.GaugeValue,
						float64(output.PowerOnTime.Hours*60*60),
						output.Device.Name,
					)
					return nil
				},
			},
			&infoMetric{
				PromDesc: prometheus.NewDesc(
					"smart_device_power_cycle",
					"Number of power cycles of the device",
					[]string{"device"},
					nil,
				),
				UpdateFunc: func(metrics chan<- prometheus.Metric, output smartctl.InfoAllOutput, desc *prometheus.Desc) error {
					metrics <- prometheus.MustNewConstMetric(
						desc,
						prometheus.CounterValue,
						float64(output.PowerCycleCount),
						output.Device.Name,
					)
					return nil
				},
			},
			&infoMetric{
				PromDesc: prometheus.NewDesc(
					"smart_device_temperature",
					"Current temperature of the device",
					[]string{"device"},
					nil,
				),
				UpdateFunc: func(metrics chan<- prometheus.Metric, output smartctl.InfoAllOutput, desc *prometheus.Desc) error {
					metrics <- prometheus.MustNewConstMetric(
						desc,
						prometheus.CounterValue,
						float64(output.Temperature.Current),
						output.Device.Name,
					)
					return nil
				},
			},
			&infoMetric{
				PromDesc: prometheus.NewDesc(
					"smart_device_raw_read_error_rate",
					"Raw read error rate of the device",
					[]string{"device"},
					nil,
				),
				UpdateFunc: func(metrics chan<- prometheus.Metric, output smartctl.InfoAllOutput, desc *prometheus.Desc) error {
					for _, e := range output.AtaSmartAttributes.Table {
						if e.Name == "Raw_Read_Error_Rate" {
							metrics <- prometheus.MustNewConstMetric(
								desc,
								prometheus.GaugeValue,
								float64(e.Value),
								output.Device.Name,
							)
							return nil
						}
					}

					return nil
				},
			},
			&infoMetric{
				PromDesc: prometheus.NewDesc(
					"smart_device_seek_error_rate",
					"Seek error rate of the device",
					[]string{"device"},
					nil,
				),
				UpdateFunc: func(metrics chan<- prometheus.Metric, output smartctl.InfoAllOutput, desc *prometheus.Desc) error {
					for _, e := range output.AtaSmartAttributes.Table {
						if e.Name == "Seek_Error_Rate" {
							metrics <- prometheus.MustNewConstMetric(
								desc,
								prometheus.GaugeValue,
								float64(e.Value),
								output.Device.Name,
							)
							return nil
						}
					}

					return nil
				},
			},
			&infoMetric{
				PromDesc: prometheus.NewDesc(
					"smart_device_reallocated_sector_count",
					"Reallocated sector count of the device",
					[]string{"device"},
					nil,
				),
				UpdateFunc: func(metrics chan<- prometheus.Metric, output smartctl.InfoAllOutput, desc *prometheus.Desc) error {
					for _, e := range output.AtaSmartAttributes.Table {
						if e.Name == "Reallocated_Sector_Ct" {
							metrics <- prometheus.MustNewConstMetric(
								desc,
								prometheus.GaugeValue,
								float64(e.Value),
								output.Device.Name,
							)
							return nil
						}
					}

					return nil
				},
			},
			&infoMetric{
				PromDesc: prometheus.NewDesc(
					"smart_device_spin_up_time",
					"Spin up time of the device",
					[]string{"device"},
					nil,
				),
				UpdateFunc: func(metrics chan<- prometheus.Metric, output smartctl.InfoAllOutput, desc *prometheus.Desc) error {
					for _, e := range output.AtaSmartAttributes.Table {
						if e.Name == "Spin_Up_Time" {
							metrics <- prometheus.MustNewConstMetric(
								desc,
								prometheus.GaugeValue,
								float64(e.Value),
								output.Device.Name,
							)
							return nil
						}
					}

					return nil
				},
			},
		},
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
