package routes

import (
	"net/http"
	"strconv"

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

func GetConstellationByID(c *gin.Context) {
	constellationID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid constellation ID"})
		return
	}

	constellation, err := db.GetConstellation(constellationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if constellation == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Constellation not found"})
		return
	}

	c.JSON(http.StatusOK, constellation)
}

func GetConstellationsByRegion(c *gin.Context) {
	regionID, err := strconv.Atoi(c.Param("regionID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid region ID"})
		return
	}

	constellations, err := db.GetConstellationsByRegionID(regionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, constellations)
}
