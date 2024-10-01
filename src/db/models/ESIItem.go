package models

type ESIItem struct {
	TypeID         int     `gorm:"primaryKey" json:"type_id"`
	GroupID        int     `gorm:"index" json:"group_id"`
	Name           string  `gorm:"type:text" json:"name"`
	Description    string  `gorm:"type:text" json:"description"`
	Mass           float64 `json:"mass"`
	Volume         float64 `json:"volume"`
	Capacity       float64 `json:"capacity"`
	PortionSize    int     `json:"portion_size"`
	PackagedVolume float64 `json:"packaged_volume"`
	Published      bool    `json:"published"`
	Radius         float64 `json:"radius"`
}

type ZKillboardItem struct {
	Flag              int              `json:"flag"`
	ItemTypeID        int              `json:"item_type_id"`
	QuantityDestroyed *int64           `json:"quantity_destroyed,omitempty"`
	QuantityDropped   *int64           `json:"quantity_dropped,omitempty"`
	Singleton         int              `json:"singleton"`
	Items             []ZKillboardItem `gorm:"type:jsonb" json:"items,omitempty"`
}
