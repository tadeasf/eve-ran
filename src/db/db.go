package db

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/tadeasf/eve-ran/src/db/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func InsertCharacter(character *models.Character) error {
	return DB.Create(character).Error
}

func GetCharacterByID(id int64) (*models.Character, error) {
	var character models.Character
	err := DB.First(&character, id).Error
	return &character, err
}

func InsertKill(kill *models.Kill) error {
	result := DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "killmail_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"character_id", "kill_time", "solar_system_id", "location_id", "hash", "fitted_value", "dropped_value", "destroyed_value", "total_value", "points", "npc", "solo", "awox", "victim_alliance_id", "victim_character_id", "victim_corporation_id", "victim_faction_id", "victim_damage_taken", "victim_ship_type_id", "victim_items", "victim_position", "attackers"}),
	}).Create(kill)

	if result.Error != nil {
		return fmt.Errorf("error upserting kill: %v", result.Error)
	}
	return nil
}

func GetKillByID(id int64) (*models.Kill, error) {
	var kill models.Kill
	err := DB.First(&kill, id).Error
	return &kill, err
}

func GetLastKillTimeForCharacter(characterID int64) (time.Time, error) {
	var lastKill struct {
		KillTime time.Time
	}

	result := DB.Table("kills").
		Where("character_id = ?", characterID).
		Order("kill_time DESC").
		Limit(1).
		Select("kill_time").
		Scan(&lastKill)

	if result.Error != nil {
		return time.Time{}, result.Error
	}

	if result.RowsAffected == 0 {
		return time.Time{}, nil
	}

	return lastKill.KillTime, nil
}

func UpsertRegion(region *models.Region) error {
	constellationsJSON, err := json.Marshal(region.Constellations)
	if err != nil {
		return err
	}

	return DB.Exec(`
        INSERT INTO regions (region_id, name, description, constellations)
        VALUES (?, ?, ?, ?)
        ON CONFLICT (region_id) DO UPDATE
        SET name = EXCLUDED.name,
            description = EXCLUDED.description,
            constellations = EXCLUDED.constellations
    `, region.RegionID, region.Name, region.Description, constellationsJSON).Error
}

func GetAllRegions() ([]models.Region, error) {
	var regions []models.Region
	err := DB.Find(&regions).Error
	if err != nil {
		return nil, err
	}
	return regions, nil
}

func UpsertSystem(system *models.System) error {
	planetsJSON, err := json.Marshal(system.Planets)
	if err != nil {
		return err
	}

	stargatesJSON, err := json.Marshal(system.Stargates)
	if err != nil {
		return err
	}

	stationsJSON, err := json.Marshal(system.Stations)
	if err != nil {
		return err
	}

	positionJSON, err := json.Marshal(system.Position)
	if err != nil {
		return err
	}

	return DB.Exec(`
        INSERT INTO systems (system_id, constellation_id, name, security_class, security_status, star_id, planets, stargates, stations, position)
        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
        ON CONFLICT (system_id) DO UPDATE
        SET constellation_id = EXCLUDED.constellation_id,
            name = EXCLUDED.name,
            security_class = EXCLUDED.security_class,
            security_status = EXCLUDED.security_status,
            star_id = EXCLUDED.star_id,
            planets = EXCLUDED.planets,
            stargates = EXCLUDED.stargates,
            stations = EXCLUDED.stations,
            position = EXCLUDED.position
    `, system.SystemID, system.ConstellationID, system.Name, system.SecurityClass, system.SecurityStatus, system.StarID, planetsJSON, stargatesJSON, stationsJSON, positionJSON).Error
}

func GetAllSystems() ([]models.System, error) {
	var systems []models.System
	err := DB.Find(&systems).Error
	return systems, err
}

func UpsertConstellation(constellation *models.Constellation) error {
	systemsJSON, err := json.Marshal(constellation.Systems)
	if err != nil {
		return err
	}

	return DB.Exec(`
        INSERT INTO constellations (constellation_id, name, region_id, systems, position)
        VALUES (?, ?, ?, ?, ?)
        ON CONFLICT (constellation_id) DO UPDATE
        SET name = EXCLUDED.name,
            region_id = EXCLUDED.region_id,
            systems = EXCLUDED.systems,
            position = EXCLUDED.position
    `, constellation.ConstellationID, constellation.Name, constellation.RegionID, systemsJSON, constellation.Position).Error
}

func GetAllConstellations() ([]models.Constellation, error) {
	var constellations []models.Constellation
	err := DB.Find(&constellations).Error
	return constellations, err
}

func UpsertESIItem(item *models.ESIItem) error {
	return DB.Save(item).Error
}

func GetAllESIItems() ([]models.ESIItem, error) {
	var items []models.ESIItem
	err := DB.Find(&items).Error
	return items, err
}

func GetESIItemByTypeID(typeID int) (*models.ESIItem, error) {
	var item models.ESIItem
	err := DB.First(&item, typeID).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &item, err
}

func GetKillsForCharacter(characterID int64, page, pageSize int) ([]models.Kill, error) {
	var kills []models.Kill
	offset := (page - 1) * pageSize
	err := DB.Where("character_id = ?", characterID).Offset(offset).Limit(pageSize).Find(&kills).Error
	return kills, err
}

func GetTotalKillsForCharacter(characterID int64) (int64, error) {
	var count int64
	err := DB.Model(&models.Kill{}).Where("character_id = ?", characterID).Count(&count).Error
	return count, err
}

func GetAllCharacters() ([]models.Character, error) {
	var characters []models.Character
	err := DB.Find(&characters).Error
	return characters, err
}

func GetKillByKillmailID(killmailID int64) (*models.Kill, error) {
	var kill models.Kill
	err := DB.First(&kill, killmailID).Error
	return &kill, err
}

func UpsertKill(kill *models.Kill) error {
	result := DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "killmail_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"character_id", "kill_time", "solar_system_id", "location_id", "hash", "fitted_value", "dropped_value", "destroyed_value", "total_value", "points", "npc", "solo", "awox", "victim_alliance_id", "victim_character_id", "victim_corporation_id", "victim_faction_id", "victim_damage_taken", "victim_ship_type_id", "victim_items", "victim_position", "attackers"}),
	}).Create(kill)

	if result.Error != nil {
		return fmt.Errorf("error upserting kill: %v", result.Error)
	}

	return nil
}

func GetAllKills() ([]models.Kill, error) {
	var kills []models.Kill
	err := DB.Find(&kills).Error
	return kills, err
}
