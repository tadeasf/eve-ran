package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tadeasf/eve-ran/src/db"
	"github.com/tadeasf/eve-ran/src/services"
)

func FetchAndStoreSystems(c *gin.Context) {
	systems, err := services.FetchAllSystems(30) // Use 30 concurrent requests
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for _, system := range systems {
		err = db.UpsertSystem(system)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Systems fetched and stored successfully", "count": len(systems)})
}

func GetAllSystems(c *gin.Context) {
	systems, err := db.GetAllSystems()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, systems)
}
