package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tadeasf/eve-ran/src/db"
	"github.com/tadeasf/eve-ran/src/services"
)

func FetchAndStoreRegions(c *gin.Context) {
	regions, err := services.FetchAllRegions(10) // Use 10 concurrent requests
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for _, region := range regions {
		err = db.UpsertRegion(region)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Regions fetched and stored successfully", "count": len(regions)})
}

// GetAllRegions retrieves all regions from the database
// @Summary Get all regions
// @Description Fetch all regions from the database
// @Tags regions
// @Accept json
// @Produce json
// @Success 200 {array} models.Region
// @Failure 500 {object} models.ErrorResponse
// @Router /regions [get]
func GetAllRegions(c *gin.Context) {
	regions, err := db.GetAllRegions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, regions)
}
