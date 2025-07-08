package service

import (
	"github.com/ethrgeist/brickstream-exporter/models"
	"github.com/ethrgeist/brickstream-exporter/repository"
	"github.com/rs/zerolog"
)

type BrickstreamService interface {
	SaveMetrics(m *models.MetricsV5) error
}

type brickstreamService struct {
	sr     repository.SiteRepository
	dr     repository.DeviceRepository
	cr     repository.CounterRepository
	logger zerolog.Logger
}

func NewBrickstreamService(sr repository.SiteRepository, dr repository.DeviceRepository, cr repository.CounterRepository, logger zerolog.Logger) BrickstreamService {
	return &brickstreamService{
		sr:     sr,
		dr:     dr,
		cr:     cr,
		logger: logger,
	}
}

func (b brickstreamService) SaveMetrics(m *models.MetricsV5) error {
	b.logger.Info().Msg("Saving metrics")
	site := &models.Site{
		SiteID:     m.SiteID,
		SiteName:   m.SiteName,
		DivisionID: m.DivisionID,
	}
	err := b.sr.Upsert(site)
	if err != nil {
		b.logger.Error().Err(err).Msg("Failed to get site by SiteID")
		return err
	}
	b.logger.Info().Str("ID", site.ID).Msg("Site upserted successfully")

	device := &models.Device{
		Name:         m.DeviceName,
		MacAddress:   m.Properties.MacAddress,
		IPAddress:    m.Properties.IPAddress,
		HostName:     m.Properties.HostName,
		HTTPPort:     m.Properties.HTTPPort,
		HTTPSPort:    m.Properties.HTTPSPort,
		Timezone:     m.Properties.Timezone,
		TimezoneName: m.Properties.TimezoneName,
		DST:          m.Properties.DST,
		HwPlatform:   m.Properties.HwPlatform,
		SerialNumber: m.Properties.SerialNumber,
		DeviceType:   m.Properties.DeviceType,
		SwRelease:    m.Properties.SwRelease,
		LastTransmit: m.Properties.TransmitTimeUTC,
		SiteID:       site.ID,
		Site:         site,
	}

	err = b.dr.Upsert(device)
	if err != nil {
		b.logger.Error().Err(err).Msg("Failed to upsert device")
		return err
	}

	b.logger.Info().Str("ID", site.ID).Msg("Device upserted successfully")

	// only take the first report and object for simplicity, must be expanded for multiple reports/objects
	counter := &models.Counter{
		SiteID:    site.ID,
		Site:      site,
		DeviceID:  device.ID,
		Device:    device,
		StartTime: m.ReportData.Reports[0].Objects[0].Counts[0].StartTimeUTC,
		EndTime:   m.ReportData.Reports[0].Objects[0].Counts[0].EndTimeUTC,
		Enters:    m.ReportData.Reports[0].Objects[0].Counts[0].Enters,
		Exits:     m.ReportData.Reports[0].Objects[0].Counts[0].Exits,
		Status:    m.ReportData.Reports[0].Objects[0].Counts[0].Status,
	}

	err = b.cr.Create(counter)
	if err != nil {
		b.logger.Error().Err(err).Msg("Failed to save counter")
		return err
	}

	return nil
}
