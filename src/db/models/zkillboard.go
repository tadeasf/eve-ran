package models

import (
	"time"
)

type Character struct {
	ID int64 `gorm:"primaryKey" json:"id"`
}

type Kill struct {
	KillmailID     int64      `gorm:"primaryKey" json:"killmail_id"`
	CharacterID    int64      `gorm:"index" json:"character_id"`
	KillTime       time.Time  `json:"killmail_time"`
	SolarSystemID  int        `gorm:"index" json:"solar_system_id"`
	LocationID     int64      `json:"locationID"`
	Hash           string     `gorm:"type:text" json:"hash"`
	FittedValue    float64    `json:"fittedValue"`
	DroppedValue   float64    `json:"droppedValue"`
	DestroyedValue float64    `json:"destroyedValue"`
	TotalValue     float64    `json:"totalValue"`
	Points         int        `json:"points"`
	NPC            bool       `json:"npc"`
	Solo           bool       `json:"solo"`
	Awox           bool       `json:"awox"`
	Victim         Victim     `gorm:"type:jsonb" json:"victim"`
	Attackers      []Attacker `gorm:"type:jsonb" json:"attackers"`
}

type Victim struct {
	AllianceID    *int      `json:"alliance_id,omitempty"`
	CharacterID   *int      `json:"character_id,omitempty"`
	CorporationID *int      `json:"corporation_id,omitempty"`
	FactionID     *int      `json:"faction_id,omitempty"`
	DamageTaken   int       `json:"damage_taken"`
	ShipTypeID    int       `json:"ship_type_id"`
	Items         []Item    `json:"items,omitempty"`
	Position      *Position `json:"position,omitempty"`
}

type Attacker struct {
	AllianceID     *int    `json:"alliance_id,omitempty"`
	CharacterID    *int    `json:"character_id,omitempty"`
	CorporationID  *int    `json:"corporation_id,omitempty"`
	FactionID      *int    `json:"faction_id,omitempty"`
	DamageDone     int     `json:"damage_done"`
	FinalBlow      bool    `json:"final_blow"`
	SecurityStatus float64 `json:"security_status"`
	ShipTypeID     *int    `json:"ship_type_id,omitempty"`
	WeaponTypeID   *int    `json:"weapon_type_id,omitempty"`
}
