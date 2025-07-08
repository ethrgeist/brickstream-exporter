package service

import (
	"github.com/ethrgeist/brickstream-exporter/repository"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog"
	"strconv"
)

type MetricsService interface {
	UpdateSiteMetrics()
}

type metricsService struct {
	sr repository.SiteRepository
	dr repository.DeviceRepository
	cr repository.CounterRepository

	sites        *prometheus.GaugeVec
	devices      *prometheus.GaugeVec
	counterEnter *prometheus.GaugeVec
	counterExit  *prometheus.GaugeVec

	log zerolog.Logger
}

func NewMetricsService(sr repository.SiteRepository, dr repository.DeviceRepository, cr repository.CounterRepository, log zerolog.Logger) MetricsService {
	sites := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "brickstream_site",
			Help: "Status of each site, value is always 1",
		},
		[]string{"site_id", "site_name", "division_id"},
	)

	// realistically, this behaves as a counter, since last transit time only increases,
	// but it's a gauge here to just use set to update the last time we saw the device
	devices := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "brickstream_device",
			Help: "Devices known to the exporter, gauge is time of last update",
		},
		[]string{
			"mac_address",
			"ip_address",
			"host_name",
			"http_port",
			"https_port",
			"timezone",
			"timezone_name",
			"dst",
			"hw_platform",
			"serial_number",
			"device_type",
			"sw_release",
		},
	)

	counterEnter := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "brickstream_counter_enters",
			Help: "Enters during the last period",
		},
		[]string{
			"site_id",
			"site_name",
			"device_hostname",
			"start_time",
			"end_time",
			"status",
			"device_name",
		},
	)

	counterExit := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "brickstream_counter_exits",
			Help: "Exists during the last period",
		},
		[]string{
			"site_id",
			"site_name",
			"device_hostname",
			"start_time",
			"end_time",
			"status",
			"device_name",
		},
	)

	prometheus.MustRegister(sites)
	prometheus.MustRegister(devices)
	prometheus.MustRegister(counterEnter)
	prometheus.MustRegister(counterExit)

	return &metricsService{
		sr:           sr,
		dr:           dr,
		cr:           cr,
		sites:        sites,
		devices:      devices,
		counterEnter: counterEnter,
		counterExit:  counterExit,
		log:          log,
	}
}

func (s *metricsService) UpdateSiteMetrics() {
	sites, err := s.sr.All()
	if err != nil {
		s.log.Error().Err(err).Msg("Failed to fetch sites for metrics update")
		return
	}
	for _, site := range sites {
		s.sites.WithLabelValues(
			site.SiteID,
			site.SiteName,
			site.DivisionID,
		).Set(1)
	}

	devices, err := s.dr.All()
	if err != nil {
		s.log.Error().Err(err).Msg("Failed to fetch devices for metrics update")
		return
	}

	for _, device := range devices {
		s.devices.WithLabelValues(
			device.MacAddress,
			device.IPAddress,
			device.HostName,
			strconv.Itoa(device.HTTPPort),
			strconv.Itoa(device.HTTPSPort),
			strconv.Itoa(device.Timezone),
			device.TimezoneName,
			strconv.Itoa(device.DST),
			device.HwPlatform,
			device.SerialNumber,
			strconv.Itoa(device.DeviceType),
			device.SwRelease,
		).Set(float64(device.LastTransmit.Unix()))

		counter, err := s.cr.GetLatestByDevice(device.ID)
		if err != nil {
			s.log.Error().Err(err).Str("device_id", device.ID).Msg("Failed to fetch latest counter for device")
			continue
		}

		if counter == nil {
			s.log.Warn().Str("device_id", device.ID).Msg("No counters found for device, skipping")
			continue
		}

		s.counterExit.WithLabelValues(
			counter.Site.SiteID,
			counter.Site.SiteName,
			counter.Device.HostName,
			counter.StartTime.String(),
			counter.EndTime.String(),
			strconv.Itoa(counter.Status),
			counter.Device.Name,
		).Set(float64(counter.Exits))

		s.counterEnter.WithLabelValues(
			counter.Site.SiteID,
			counter.Site.SiteName,
			counter.Device.HostName,
			counter.StartTime.String(),
			counter.EndTime.String(),
			strconv.Itoa(counter.Status),
			counter.Device.Name,
		).Set(float64(counter.Enters))

	}
}
