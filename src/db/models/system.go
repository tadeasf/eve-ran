package models

type System struct {
	SystemID        int      `gorm:"primaryKey" json:"system_id"`
	ConstellationID int      `gorm:"index" json:"constellation_id"`
	Name            string   `gorm:"type:text" json:"name"`
	SecurityClass   string   `gorm:"type:text" json:"security_class"`
	SecurityStatus  float64  `json:"security_status"`
	StarID          int      `json:"star_id"`
	Planets         IntArray `gorm:"type:jsonb" json:"planets"`
	Stargates       IntArray `gorm:"type:jsonb" json:"stargates"`
	Stations        IntArray `gorm:"type:jsonb" json:"stations"`
	Position        Position `gorm:"type:jsonb" json:"position"`
}
