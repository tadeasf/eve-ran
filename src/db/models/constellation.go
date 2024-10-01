package models

type Constellation struct {
	ConstellationID int      `gorm:"primaryKey" json:"constellation_id"`
	Name            string   `gorm:"type:text" json:"name"`
	RegionID        int      `gorm:"index" json:"region_id"`
	Systems         IntArray `gorm:"type:jsonb" json:"systems"`
	Position        Position `gorm:"type:jsonb" json:"position"`
}
