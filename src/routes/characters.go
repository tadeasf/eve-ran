package routes

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tadeasf/eve-ran/src/db"
)

// GetAllCharacters retrieves all characters from the database
// @Summary Get all characters
// @Description Fetch all characters from the database
// @Tags characters
// @Accept json
// @Produce json
// @Success 200 {array} models.Character
// @Failure 500 {object} models.ErrorResponse
// @Router /characters [get]
func GetAllCharacters(c *gin.Context) {
	characters, err := db.GetAllCharacters()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, characters)
}

// GetAllKills retrieves all kills from the database
// @Summary Get all kills
// @Description Fetch all kills from the database
// @Tags kills
// @Accept json
// @Produce json
// @Success 200 {array} models.Kill
// @Failure 500 {object} models.ErrorResponse
// @Router /kills [get]
func GetAllKills(c *gin.Context) {
	kills, err := db.GetAllKills()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, kills)
}

// GetAllCharacterStats retrieves stats for all characters with filters
// @Summary Get all character stats
// @Description Fetch stats for all characters from the database with optional filters
// @Tags characters
// @Accept json
// @Produce json
// @Param regionID query []int false "Region IDs"
// @Param startDate query string false "Start date (YYYY-MM-DD)"
// @Param endDate query string false "End date (YYYY-MM-DD)"
// @Success 200 {array} db.CharacterStats
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /characters/stats [get]
func GetAllCharacterStats(c *gin.Context) {
	regionIDs := c.QueryArray("regionID")
	startDate := c.Query("startDate")
	endDate := c.Query("endDate")

	// Convert regionIDs from string to int
	var regionIDInts []int64
	for _, id := range regionIDs {
		intID, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid region ID"})
			return
		}
		regionIDInts = append(regionIDInts, intID)
	}

	// Parse dates
	var startTime, endTime time.Time
	var err error
	if startDate != "" {
		startTime, err = time.Parse("2006-01-02", startDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start date format"})
			return
		}
	}
	if endDate != "" {
		endTime, err = time.Parse("2006-01-02", endDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end date format"})
			return
		}
	}

	stats, err := db.GetCharacterStats(startTime, endTime, 0, regionIDInts...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, stats)
}
