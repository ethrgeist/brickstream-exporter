package models

import "gorm.io/gorm"

type Site struct {
	ID string `gorm:"primaryKey"`

	SiteID     string
	SiteName   string
	DeviceID   string
	DeviceName string
	DivisionID string

	gorm.Model
}
