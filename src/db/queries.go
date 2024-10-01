package db

import (
	"time"

	"github.com/tadeasf/eve-ran/src/db/models"
)

func GetConstellation(id int) (*models.Constellation, error) {
	var constellation models.Constellation
	result := DB.First(&constellation, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &constellation, nil
}

func GetConstellationsByRegionID(regionID int) ([]models.Constellation, error) {
	var constellations []models.Constellation
	result := DB.Where("region_id = ?", regionID).Find(&constellations)
	if result.Error != nil {
		return nil, result.Error
	}
	return constellations, nil
}

func GetSystem(id int) (*models.System, error) {
	var system models.System
	result := DB.First(&system, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &system, nil
}

func GetSystemsByRegionID(regionID int) ([]models.System, error) {
	var systems []models.System
	result := DB.Where("region_id = ?", regionID).Find(&systems)
	if result.Error != nil {
		return nil, result.Error
	}
	return systems, nil
}

func GetKillsForCharacterWithFilters(characterID int64, page, pageSize, regionID int, startDate, endDate string) ([]models.Kill, error) {
	var kills []models.Kill
	query := DB.Where("character_id = ?", characterID)

	if regionID != 0 {
		query = query.Where("solar_system_id IN (SELECT system_id FROM systems WHERE region_id = ?)", regionID)
	}

	if startDate != "" {
		startTime, _ := time.Parse("2006-01-02", startDate)
		query = query.Where("kill_time >= ?", startTime)
	}

	if endDate != "" {
		endTime, _ := time.Parse("2006-01-02", endDate)
		query = query.Where("kill_time <= ?", endTime)
	}

	result := query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&kills)
	if result.Error != nil {
		return nil, result.Error
	}
	return kills, nil
}

func GetTotalKillsForCharacterWithFilters(characterID int64, regionID int, startDate, endDate string) (int64, error) {
	var count int64
	query := DB.Model(&models.Kill{}).Where("character_id = ?", characterID)

	if regionID != 0 {
		query = query.Where("solar_system_id IN (SELECT system_id FROM systems WHERE region_id = ?)", regionID)
	}

	if startDate != "" {
		startTime, _ := time.Parse("2006-01-02", startDate)
		query = query.Where("kill_time >= ?", startTime)
	}

	if endDate != "" {
		endTime, _ := time.Parse("2006-01-02", endDate)
		query = query.Where("kill_time <= ?", endTime)
	}

	result := query.Count(&count)
	if result.Error != nil {
		return 0, result.Error
	}
	return count, nil
}

func GetKillsByRegion(regionID, page, pageSize int, startDate, endDate string) ([]models.Kill, error) {
	var kills []models.Kill
	query := DB.Where("solar_system_id IN (SELECT system_id FROM systems WHERE region_id = ?)", regionID)

	if startDate != "" {
		startTime, _ := time.Parse("2006-01-02", startDate)
		query = query.Where("kill_time >= ?", startTime)
	}

	if endDate != "" {
		endTime, _ := time.Parse("2006-01-02", endDate)
		query = query.Where("kill_time <= ?", endTime)
	}

	result := query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&kills)
	if result.Error != nil {
		return nil, result.Error
	}
	return kills, nil
}

func GetTotalKillsByRegion(regionID int, startDate, endDate string) (int64, error) {
	var count int64
	query := DB.Model(&models.Kill{}).Where("solar_system_id IN (SELECT system_id FROM systems WHERE region_id = ?)", regionID)

	if startDate != "" {
		startTime, _ := time.Parse("2006-01-02", startDate)
		query = query.Where("kill_time >= ?", startTime)
	}

	if endDate != "" {
		endTime, _ := time.Parse("2006-01-02", endDate)
		query = query.Where("kill_time <= ?", endTime)
	}

	result := query.Count(&count)
	if result.Error != nil {
		return 0, result.Error
	}
	return count, nil
}
