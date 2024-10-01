package models

import "time"

type Character struct {
	ID int64 `json:"id" db:"id,bigint,PRIMARY KEY"`
}

type Kill struct {
	KillmailID     int64      `json:"killmail_id" db:"killmail_id,bigint,PRIMARY KEY"`
	CharacterID    int64      `json:"character_id" db:"character_id,bigint,REFERENCES characters(id)"`
	KillTime       time.Time  `json:"killmail_time" db:"killmail_time,timestamp"`
	SolarSystemID  int        `json:"solar_system_id" db:"solar_system_id,integer"`
	LocationID     int64      `json:"locationID" db:"location_id,bigint"`
	Hash           string     `json:"hash" db:"hash,text"`
	FittedValue    float64    `json:"fittedValue" db:"fitted_value,numeric"`
	DroppedValue   float64    `json:"droppedValue" db:"dropped_value,numeric"`
	DestroyedValue float64    `json:"destroyedValue" db:"destroyed_value,numeric"`
	TotalValue     float64    `json:"totalValue" db:"total_value,numeric"`
	Points         int        `json:"points" db:"points,integer"`
	NPC            bool       `json:"npc" db:"npc,boolean"`
	Solo           bool       `json:"solo" db:"solo,boolean"`
	Awox           bool       `json:"awox" db:"awox,boolean"`
	Victim         Victim     `json:"victim" db:"victim,jsonb"`
	Attackers      []Attacker `json:"attackers" db:"attackers,jsonb"`
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

type Item struct {
	Flag              int    `json:"flag"`
	ItemTypeID        int    `json:"item_type_id"`
	QuantityDestroyed *int64 `json:"quantity_destroyed,omitempty"`
	QuantityDropped   *int64 `json:"quantity_dropped,omitempty"`
	Singleton         int    `json:"singleton"`
	Items             []Item `json:"items,omitempty"`
}

type Position struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
}
