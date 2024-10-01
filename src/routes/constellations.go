package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tadeasf/eve-ran/src/db"
	"github.com/tadeasf/eve-ran/src/services"
)

func FetchAndStoreConstellations(c *gin.Context) {
	constellations, err := services.FetchAllConstellations(20) // Use 20 concurrent requests
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for _, constellation := range constellations {
		err = db.UpsertConstellation(constellation)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Constellations fetched and stored successfully", "count": len(constellations)})
}

func GetAllConstellations(c *gin.Context) {
	constellations, err := db.GetAllConstellations()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, constellations)
}
