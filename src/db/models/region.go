package models

import (
	"database/sql/driver"
	"encoding/json"
)

type Region struct {
	RegionID       int      `gorm:"primaryKey" json:"region_id"`
	Name           string   `gorm:"type:text" json:"name"`
	Description    string   `gorm:"type:text" json:"description"`
	Constellations IntArray `gorm:"type:jsonb" json:"constellations"`
}

func (a IntArray) Value() (driver.Value, error) {
	return json.Marshal(a)
}
