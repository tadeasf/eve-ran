package models

type Region struct {
	RegionID       int    `gorm:"primaryKey" json:"region_id"`
	Name           string `gorm:"type:text" json:"name"`
	Description    string `gorm:"type:text" json:"description"`
	Constellations []int  `gorm:"type:jsonb" json:"constellations"`
}
