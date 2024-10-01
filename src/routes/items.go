package routes

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tadeasf/eve-ran/src/db"
	"github.com/tadeasf/eve-ran/src/services"
)

func FetchAndStoreItems(c *gin.Context) {
	items, err := services.FetchAllItems(50) // Use 50 concurrent requests
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for _, item := range items {
		err = db.UpsertESIItem(item)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Items fetched and stored successfully", "count": len(items)})
}

func GetAllItems(c *gin.Context) {
	items, err := db.GetAllESIItems()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, items)
}

func GetItemByTypeID(c *gin.Context) {
	typeID := c.Param("typeID")
	itemTypeID, err := strconv.Atoi(typeID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid type ID"})
		return
	}

	item, err := db.GetESIItemByTypeID(itemTypeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if item == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
		return
	}

	c.JSON(http.StatusOK, item)
}
