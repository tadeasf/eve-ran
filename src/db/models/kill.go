package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

type Character struct {
	ID int64 `gorm:"primaryKey" json:"id"`
}

type Kill struct {
	KillmailID     int64         `json:"killmail_id" gorm:"primaryKey"`
	CharacterID    int64         `json:"character_id"`
	KillTime       time.Time     `json:"killmail_time"`
	SolarSystemID  int           `json:"solar_system_id"`
	LocationID     int64         `json:"locationID"`
	Hash           string        `json:"hash"`
	FittedValue    float64       `json:"fitted_value"`
	DroppedValue   float64       `json:"dropped_value"`
	DestroyedValue float64       `json:"destroyed_value"`
	TotalValue     float64       `json:"total_value"`
	Points         int           `json:"points"`
	NPC            bool          `json:"npc"`
	Solo           bool          `json:"solo"`
	Awox           bool          `json:"awox"`
	Victim         Victim        `json:"victim" gorm:"embedded;embeddedPrefix:victim_"`
	Attackers      AttackersJSON `json:"attackers" gorm:"type:jsonb"`
}

type AttackersJSON []Attacker

func (a AttackersJSON) Value() (driver.Value, error) {
	return json.Marshal(a)
}

func (a *AttackersJSON) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}

	return json.Unmarshal(bytes, &a)
}

type Victim struct {
	AllianceID    *int      `json:"alliance_id,omitempty"`
	CharacterID   *int      `json:"character_id,omitempty"`
	CorporationID *int      `json:"corporation_id,omitempty"`
	FactionID     *int      `json:"faction_id,omitempty"`
	DamageTaken   int       `json:"damage_taken"`
	ShipTypeID    int       `json:"ship_type_id"`
	Items         ItemArray `json:"items" gorm:"type:jsonb"`
	Position      *Position `json:"position" gorm:"type:jsonb"`
}

type Attacker struct {
	AllianceID     *int    `json:"alliance_id,omitempty"`
	CharacterID    *int    `json:"character_id,omitempty"`
	CorporationID  *int    `json:"corporation_id,omitempty"`
	FactionID      *int    `json:"faction_id,omitempty"`
	DamageDone     int     `json:"damage_done"`
	FinalBlow      bool    `json:"final_blow"`
	SecurityStatus float64 `json:"security_status"`
	ShipTypeID     int     `json:"ship_type_id"`
	WeaponTypeID   int     `json:"weapon_type_id"`
}

type Item struct {
	ItemTypeID        int  `json:"item_type_id"`
	Singleton         int  `json:"singleton"`
	QuantityDropped   *int `json:"quantity_dropped,omitempty"`
	QuantityDestroyed *int `json:"quantity_destroyed,omitempty"`
	Flag              int  `json:"flag"`
}

func (a Attacker) Value() (driver.Value, error) {
	return json.Marshal(a)
}

func (a *Attacker) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}

	return json.Unmarshal(bytes, &a)
}

func (k *Kill) MarshalJSON() ([]byte, error) {
	type Alias Kill
	return json.Marshal(&struct {
		*Alias
		Attackers json.RawMessage `json:"attackers"`
	}{
		Alias:     (*Alias)(k),
		Attackers: k.AttackersJSON(),
	})
}

func (k *Kill) UnmarshalJSON(data []byte) error {
	type Alias Kill
	aux := &struct {
		*Alias
		Attackers json.RawMessage `json:"attackers"`
	}{
		Alias: (*Alias)(k),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	return json.Unmarshal(aux.Attackers, &k.Attackers)
}

func (k *Kill) AttackersJSON() json.RawMessage {
	b, _ := json.Marshal(k.Attackers)
	return json.RawMessage(b)
}
