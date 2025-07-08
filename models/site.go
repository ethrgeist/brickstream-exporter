package models

import "gorm.io/gorm"

type Site struct {
	ID string `gorm:"primaryKey"`

	SiteID     string
	SiteName   string
	DivisionID string

	gorm.Model
}
