package repository

import (
	"errors"
	"github.com/ethrgeist/brickstream-exporter/models"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type SiteRepository interface {
	Create(site *models.Site) error
	Update(site *models.Site) error
	Upsert(site *models.Site) error
	Delete(siteID string) error
	All() ([]*models.Site, error)
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

func (sr *siteRepository) Upsert(site *models.Site) error {
	existingSite, err := sr.GetBySiteID(site.SiteID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return sr.Create(site)
		}
		if err != nil {
			sr.log.Error().Err(err).Msg("Failed to get site by SiteID")
			return err
		}
	}
	existingSite.SiteName = site.SiteName
	existingSite.DivisionID = site.DivisionID
	result := sr.db.Save(existingSite)
	if result.Error != nil {
		sr.log.Error().Err(result.Error).Msg("Failed to upsert site")
		return result.Error
	}
	*site = *existingSite
	sr.log.Debug().Str("SiteID", site.SiteID).Msg("Site upserted successfully")
	return nil
}

func (sr *siteRepository) Delete(id string) error {
	result := sr.db.Where("id = ?", id).Delete(&models.Site{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (sr *siteRepository) All() ([]*models.Site, error) {
	var sites []*models.Site
	result := sr.db.Find(&sites)
	if result.Error != nil {
		return nil, result.Error
	}
	return sites, nil
}

func (sr *siteRepository) GetBySiteID(siteID string) (*models.Site, error) {
	site := &models.Site{}
	result := sr.db.Where("site_id = ?", siteID).First(site)
	if result.Error != nil {
		return nil, result.Error
	}
	return site, nil
}
