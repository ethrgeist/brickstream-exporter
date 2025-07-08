package repository

import (
	"errors"
	"github.com/ethrgeist/brickstream-exporter/models"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type CounterRepository interface {
	Create(counter *models.Counter) error
	Update(counter *models.Counter) error
	Delete(counterID string) error
	All() ([]*models.Counter, error)
	GetLatestByDevice(deviceID string) (*models.Counter, error)
	Current() ([]*models.Counter, error)
}

type counterRepository struct {
	db  *gorm.DB
	log zerolog.Logger
}

func NewCounterRepository(db *gorm.DB, log zerolog.Logger) CounterRepository {
	return &counterRepository{
		db:  db,
		log: log,
	}
}

func (cr *counterRepository) Create(counter *models.Counter) error {
	result := cr.db.Create(counter)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (cr *counterRepository) Update(counter *models.Counter) error {
	if counter.ID == "" {
		return cr.Create(counter)
	}
	return cr.Update(counter)
}

func (cr *counterRepository) Delete(id string) error {
	result := cr.db.Where("id = ?", id).Delete(&models.Counter{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (cr *counterRepository) All() ([]*models.Counter, error) {
	var counters []*models.Counter
	result := cr.db.Find(&counters)
	if result.Error != nil {
		return nil, result.Error
	}
	return counters, nil
}

func (cr *counterRepository) GetLatestByDevice(deviceID string) (*models.Counter, error) {
	var counter models.Counter
	result := cr.db.
		Where("device_id = ?", deviceID).
		Order("created_at DESC").
		Preload("Site").
		Preload("Device").
		First(&counter)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &counter, nil
}

func (cr *counterRepository) Current() ([]*models.Counter, error) {
	var counters []*models.Counter

	subQuery := cr.db.Model(&models.Counter{}).
		Select("MAX(id)").
		Group("device_id")

	result := cr.db.
		Preload("Site").
		Preload("Device").
		Where("id IN (?)", subQuery).
		Find(&counters)

	if result.Error != nil {
		return nil, result.Error
	}
	return counters, nil
}
