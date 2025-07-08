package models

import (
	"gorm.io/gorm"
	"time"
)

type Counter struct {
	ID string `gorm:"primaryKey"`

	SiteID   string
	Site     *Site
	DeviceID string
	Device   *Device

	StartTime time.Time
	EndTime   time.Time
	Enters    int
	Exits     int
	Status    int

	gorm.Model
}
