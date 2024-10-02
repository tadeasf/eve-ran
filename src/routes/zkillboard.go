package routes

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tadeasf/eve-ran/src/db"
	"github.com/tadeasf/eve-ran/src/db/models"
	"github.com/tadeasf/eve-ran/src/jobs"
)

// AddCharacter adds a new character ID
// @Summary Add a new character ID
// @Description Add a new character ID to the database and fetch all kills
// @Tags characters
// @Accept json
// @Produce json
// @Param character body models.Character true "Character ID"
// @Success 201 {object} models.Character
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /characters [post]
func AddCharacter(c *gin.Context) {
	var character models.Character
	if err := c.ShouldBindJSON(&character); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Insert the character ID into the database
	err := db.InsertCharacter(&character)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add character"})
		return
	}

	// Trigger a full kill fetch for the new character
	go jobs.FetchAllKillsForCharacter(character.ID)

	c.JSON(http.StatusCreated, character)
}

// RemoveCharacter removes a character
// @Summary Remove a character
// @Description Remove a character from the database
// @Tags characters
// @Accept json
// @Produce json
// @Param id path int true "Character ID"
// @Success 204 "No Content"
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /characters/{id} [delete]
func RemoveCharacter(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid character ID"})
		return
	}

	err = db.DB.Delete(&models.Character{}, id).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// GetCharacterKills retrieves character kills
// @Summary Get character kills
// @Description Fetch and store kills for a character from zKillboard
// @Tags characters
// @Accept json
// @Produce json
// @Param id path int true "Character ID"
// @Param page query int false "Page number"
// @Success 200 {array} models.Kill
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /characters/{id}/kills [get]
func GetCharacterKills(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid character ID"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	kills, err := fetchKillsFromZKillboard(id, page)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = storeKills(id, kills)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, kills)
}

func fetchKillsFromZKillboard(characterID int64, page int) ([]models.Kill, error) {
	url := fmt.Sprintf("https://zkillboard.com/api/kills/characterID/%d/page/%d/", characterID, page)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var rawKills []struct {
		KillmailID    int64     `json:"killmail_id"`
		KillmailTime  time.Time `json:"killmail_time"`
		SolarSystemID int       `json:"solar_system_id"`
		ZKB           struct {
			LocationID     int64   `json:"locationID"`
			Hash           string  `json:"hash"`
			FittedValue    float64 `json:"fittedValue"`
			DroppedValue   float64 `json:"droppedValue"`
			DestroyedValue float64 `json:"destroyedValue"`
			TotalValue     float64 `json:"totalValue"`
			Points         int     `json:"points"`
			NPC            bool    `json:"npc"`
			Solo           bool    `json:"solo"`
			Awox           bool    `json:"awox"`
		} `json:"zkb"`
	}
	err = json.Unmarshal(body, &rawKills)
	if err != nil {
		return nil, err
	}

	var kills []models.Kill
	for _, rawKill := range rawKills {
		kill := models.Kill{
			KillmailID:     rawKill.KillmailID,
			CharacterID:    characterID,
			KillTime:       rawKill.KillmailTime,
			SolarSystemID:  rawKill.SolarSystemID,
			LocationID:     rawKill.ZKB.LocationID,
			Hash:           rawKill.ZKB.Hash,
			FittedValue:    rawKill.ZKB.FittedValue,
			DroppedValue:   rawKill.ZKB.DroppedValue,
			DestroyedValue: rawKill.ZKB.DestroyedValue,
			TotalValue:     rawKill.ZKB.TotalValue,
			Points:         rawKill.ZKB.Points,
			NPC:            rawKill.ZKB.NPC,
			Solo:           rawKill.ZKB.Solo,
			Awox:           rawKill.ZKB.Awox,
		}
		kills = append(kills, kill)
	}

	return kills, nil
}

func storeKills(characterID int64, kills []models.Kill) error {
	for _, kill := range kills {
		kill.CharacterID = characterID
		err := db.DB.Create(&kill).Error
		if err != nil {
			return err
		}
	}
	return nil
}

// GetCharacterKillsFromDB retrieves character kills from the database
// @Summary Get character kills from database
// @Description Fetch kills for a character from the database
// @Tags characters
// @Accept json
// @Produce json
// @Param id path int true "Character ID"
// @Param page query int false "Page number"
// @Param pageSize query int false "Page size"
// @Success 200 {object} models.PaginatedResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /characters/{id}/kills/db [get]
func GetCharacterKillsFromDB(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid character ID"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	kills, err := db.GetKillsForCharacter(id, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	totalItems, err := db.GetTotalKillsForCharacter(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	totalPages := int((totalItems + int64(pageSize) - 1) / int64(pageSize))

	response := models.PaginatedResponse{
		Data:       kills,
		Page:       page,
		PageSize:   pageSize,
		TotalItems: int(totalItems),
		TotalPages: totalPages,
	}

	c.JSON(http.StatusOK, response)
}

// GetKillsByRegion retrieves kills by region
// @Summary Get kills by region
// @Description Fetch kills for a region from the database
// @Tags kills
// @Accept json
// @Produce json
// @Param regionID path int true "Region ID"
// @Param page query int false "Page number"
// @Param pageSize query int false "Page size"
// @Param startDate query string false "Start date (YYYY-MM-DD)"
// @Param endDate query string false "End date (YYYY-MM-DD)"
// @Success 200 {object} models.PaginatedResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /kills/region/{regionID} [get]
func GetKillsByRegion(c *gin.Context) {
	regionID, err := strconv.Atoi(c.Param("regionID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid region ID"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	startDate := c.Query("startDate")
	endDate := c.Query("endDate")

	kills, totalCount, err := db.GetKillsByRegion(regionID, page, pageSize, startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	totalPages := int(math.Ceil(float64(totalCount) / float64(pageSize)))

	response := models.PaginatedResponse{
		Data:       kills,
		Page:       page,
		PageSize:   pageSize,
		TotalItems: int(totalCount),
		TotalPages: totalPages,
	}

	c.JSON(http.StatusOK, response)
}
