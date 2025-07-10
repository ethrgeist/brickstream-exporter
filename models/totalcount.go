package models

import (
	"gorm.io/gorm"
)

type TotalCount struct {
	ID string `gorm:"primaryKey"`

	DeviceID string
	Device   *Device

	Exits  int
	Enters int

	gorm.Model
}
