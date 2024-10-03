package routes

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tadeasf/eve-ran/src/db"
)

func GetCharacterKillmails(c *gin.Context) {
	characterID, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	startTime, _ := time.Parse(time.RFC3339, c.Query("start_time"))
	endTime, _ := time.Parse(time.RFC3339, c.Query("end_time"))
	systemID, _ := strconv.ParseInt(c.Query("system_id"), 10, 64)
	regionID, _ := strconv.ParseInt(c.Query("region_id"), 10, 64)

	kills, err := db.GetCharacterKillmails(characterID, startTime, endTime, systemID, regionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, kills)
}
