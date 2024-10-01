package routes

import (
	"encoding/json"
	"fmt"
	"io"
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

	_, err = db.DB.Exec("DELETE FROM characters WHERE id = $1", id)
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
		_, err := db.DB.Exec(`
			INSERT INTO kills (
				killmail_id, character_id, killmail_time, solar_system_id,
				location_id, hash, fitted_value, dropped_value, destroyed_value,
				total_value, points, npc, solo, awox
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
			ON CONFLICT (killmail_id) DO NOTHING
		`,
			kill.KillmailID, characterID, kill.KillTime, kill.SolarSystemID,
			kill.LocationID, kill.Hash, kill.FittedValue, kill.DroppedValue,
			kill.DestroyedValue, kill.TotalValue, kill.Points,
			kill.NPC, kill.Solo, kill.Awox,
		)
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

	totalPages := (totalItems + pageSize - 1) / pageSize

	response := models.PaginatedResponse{
		Data:       kills,
		Page:       page,
		PageSize:   pageSize,
		TotalItems: totalItems,
		TotalPages: totalPages,
	}

	c.JSON(http.StatusOK, response)
}
