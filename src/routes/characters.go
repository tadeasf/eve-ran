package routes

import (
	"net/http"

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
