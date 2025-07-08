package models

import (
	"gorm.io/gorm"
	"time"
)

type Device struct {
	ID string `gorm:"primaryKey"`

	Name         string
	MacAddress   string
	IPAddress    string
	HostName     string
	HTTPPort     int
	HTTPSPort    int
	Timezone     int
	TimezoneName string
	DST          int
	HwPlatform   string
	SerialNumber string
	DeviceType   int
	SwRelease    string
	LastTransmit time.Time

	SiteID string
	Site   *Site

	gorm.Model
}
