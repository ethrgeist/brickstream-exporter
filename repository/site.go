package repository

import (
	"github.com/ethrgeist/brickstream-exporter/models"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type SiteRepository interface {
	Create(site *models.Site) error
	Update(site *models.Site) error
	Delete(siteID string) error
	GetBySiteID(siteID string) (*models.Site, error)
}

type siteRepository struct {
	db  *gorm.DB
	log zerolog.Logger
}

func NewSiteRepository(db *gorm.DB, log zerolog.Logger) SiteRepository {
	return &siteRepository{
		db:  db,
		log: log,
	}
}

func (sr *siteRepository) Create(site *models.Site) error {
	result := sr.db.Create(site)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (sr *siteRepository) Update(site *models.Site) error {
	if site.ID == "" {
		return sr.Create(site)
	}
	return sr.Update(site)
}

func (sr *siteRepository) Delete(id string) error {
	result := sr.db.Where("id = ?", id).Delete(&models.Site{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (sr *siteRepository) GetBySiteID(siteID string) (*models.Site, error) {
	site := &models.Site{}
	result := sr.db.Where("site_id = ?", siteID).First(site)
	if result.Error != nil {
		return nil, result.Error
	}
	return site, nil
}
