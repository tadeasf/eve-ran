package db

import (
	"time"

	"github.com/tadeasf/eve-ran/src/db/models"
	"gorm.io/gorm"
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
	return DB.Create(kill).Error
}

func GetKillByID(id int64) (*models.Kill, error) {
	var kill models.Kill
	err := DB.First(&kill, id).Error
	return &kill, err
}

func GetLastKillTimeForCharacter(characterID int64) (time.Time, error) {
	var kill models.Kill
	err := DB.Where("character_id = ?", characterID).Order("kill_time DESC").First(&kill).Error
	return kill.KillTime, err
}

func UpsertRegion(region *models.Region) error {
	return DB.Save(region).Error
}

func GetAllRegions() ([]models.Region, error) {
	var regions []models.Region
	err := DB.Find(&regions).Error
	return regions, err
}

func UpsertSystem(system *models.System) error {
	return DB.Save(system).Error
}

func GetAllSystems() ([]models.System, error) {
	var systems []models.System
	err := DB.Find(&systems).Error
	return systems, err
}

func UpsertConstellation(constellation *models.Constellation) error {
	return DB.Save(constellation).Error
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
