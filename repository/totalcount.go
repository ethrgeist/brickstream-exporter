package repository

import (
	"errors"
	"github.com/ethrgeist/brickstream-exporter/models"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type TotalCountRepository interface {
	Upsert(*models.TotalCount) error
	GetByDeviceID(string) (*models.TotalCount, error)
}

type totalCountRepository struct {
	db  *gorm.DB
	log *zerolog.Logger
}

func NewTotalCountRepository(db *gorm.DB, log zerolog.Logger) TotalCountRepository {
	return &totalCountRepository{
		db:  db,
		log: &log,
	}
}

func (tr *totalCountRepository) Create(tc *models.TotalCount) error {
	return tr.db.Create(tc).Error
}

func (tr *totalCountRepository) Upsert(tc *models.TotalCount) error {
	existingTc, err := tr.GetByDeviceID(tc.DeviceID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return tr.Create(tc)
		}
		if err != nil {
			tr.log.Error().Err(err).Msg("Failed to get device by DeviceID")
			return err
		}
	}
	existingTc.DeviceID = tc.DeviceID
	existingTc.Exits = existingTc.Exits + tc.Exits
	existingTc.Enters = existingTc.Enters + tc.Enters

	result := tr.db.Save(existingTc)
	if result.Error != nil {
		tr.log.Error().Err(result.Error).Msg("Failed to update total count")
		return result.Error
	}

	*tc = *existingTc
	tr.log.Debug().Str("ID", existingTc.DeviceID).Msg("Device upserted successfully")
	return nil
}

func (tr *totalCountRepository) GetByDeviceID(d string) (*models.TotalCount, error) {
	var tc *models.TotalCount

	result := tr.db.Where(&models.TotalCount{DeviceID: d}).First(&tc)
	if result.Error != nil {
		return nil, result.Error
	}
	return tc, nil
}
