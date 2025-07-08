package repository

import (
	"errors"
	"github.com/ethrgeist/brickstream-exporter/models"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type DeviceRepository interface {
	Create(device *models.Device) error
	Update(device *models.Device) error
	Upsert(device *models.Device) error
	Delete(deviceID string) error
	All() ([]*models.Device, error)
	GetBySerial(deviceID string) (*models.Device, error)
}

type deviceRepository struct {
	db  *gorm.DB
	log zerolog.Logger
}

func NewDeviceRepository(db *gorm.DB, log zerolog.Logger) DeviceRepository {
	return &deviceRepository{
		db:  db,
		log: log,
	}
}

func (sr *deviceRepository) Create(device *models.Device) error {
	result := sr.db.Create(device)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (sr *deviceRepository) Update(device *models.Device) error {
	if device.ID == "" {
		return sr.Create(device)
	}
	return sr.Update(device)
}

func (sr *deviceRepository) Upsert(device *models.Device) error {
	existingDevice, err := sr.GetBySerial(device.SerialNumber)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return sr.Create(device)
		}
		if err != nil {
			sr.log.Error().Err(err).Msg("Failed to get device by DeviceID")
			return err
		}
	}
	existingDevice.IPAddress = device.IPAddress
	existingDevice.HostName = device.HostName
	existingDevice.HTTPPort = device.HTTPPort
	existingDevice.HTTPSPort = device.HTTPSPort
	existingDevice.Timezone = device.Timezone
	existingDevice.TimezoneName = device.TimezoneName
	existingDevice.DST = device.DST
	existingDevice.DeviceType = device.DeviceType
	existingDevice.SwRelease = device.SwRelease
	existingDevice.LastTransmit = device.LastTransmit

	result := sr.db.Save(existingDevice)
	if result.Error != nil {
		sr.log.Error().Err(result.Error).Msg("Failed to upsert device")
		return result.Error
	}
	*device = *existingDevice
	sr.log.Debug().Str("ID", device.ID).Msg("Device upserted successfully")
	return nil
}

func (sr *deviceRepository) Delete(id string) error {
	result := sr.db.Where("id = ?", id).Delete(&models.Device{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (sr *deviceRepository) All() ([]*models.Device, error) {
	var devices []*models.Device
	result := sr.db.Find(&devices)
	if result.Error != nil {
		return nil, result.Error
	}
	return devices, nil
}

func (sr *deviceRepository) GetBySerial(deviceID string) (*models.Device, error) {
	device := &models.Device{}
	result := sr.db.Where("serial_number = ?", deviceID).First(device)
	if result.Error != nil {
		return nil, result.Error
	}
	return device, nil
}
