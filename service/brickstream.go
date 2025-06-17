package service

import (
	"fmt"
	"github.com/ethrgeist/brickstream-exporter/models"
	"github.com/ethrgeist/brickstream-exporter/repository"
	"github.com/rs/zerolog"
)

type BrickstreamService interface {
	SaveMetrics(m *models.MetricsV5) error
}

type brickstreamService struct {
	siteRepository repository.SiteRepository
	logger         zerolog.Logger
}

func NewBrickstreamService(siteRepository repository.SiteRepository, logger zerolog.Logger) BrickstreamService {
	return &brickstreamService{
		siteRepository: siteRepository,
		logger:         logger,
	}
}

func (b brickstreamService) SaveMetrics(m *models.MetricsV5) error {
	b.logger.Info().Msg("Saving metrics")
	fmt.Printf("%+v\n", m)
	return nil
}
