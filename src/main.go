package main

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "github.com/tadeasf/eve-ran/docs"
	"github.com/tadeasf/eve-ran/src/db"
	"github.com/tadeasf/eve-ran/src/jobs"
	"github.com/tadeasf/eve-ran/src/routes"
)

// @title EVE Ran API
// @version 1.0
// @description This is the API for EVE Ran application.
// @host localhost:8080
// @BasePath /
// @schemes http https

func main() {
	db.InitDB()

	// Start the kill fetcher job
	go jobs.StartKillFetcherJob()

	// Run the type fetcher job
	go jobs.FetchAndUpdateTypes()

	r := gin.Default()

	// zKillboard routes
	r.POST("/characters", routes.AddCharacter)
	r.DELETE("/characters/:id", routes.RemoveCharacter)
	r.GET("/characters/:id/kills", routes.GetCharacterKills)
	r.GET("/characters/:id/kills/db", routes.GetCharacterKillsFromDB)

	// Region routes
	r.POST("/regions/fetch", routes.FetchAndStoreRegions)
	r.GET("/regions", routes.GetAllRegions)

	// System routes
	r.POST("/systems/fetch", routes.FetchAndStoreSystems)
	r.GET("/systems", routes.GetAllSystems)
	r.GET("/systems/:id", routes.GetSystemByID)
	r.GET("/systems/region/:regionID", routes.GetSystemsByRegion)

	// Constellation routes
	r.POST("/constellations/fetch", routes.FetchAndStoreConstellations)
	r.GET("/constellations", routes.GetAllConstellations)
	r.GET("/constellations/:id", routes.GetConstellationByID)
	r.GET("/constellations/region/:regionID", routes.GetConstellationsByRegion)

	// Item routes
	r.POST("/items/fetch", routes.FetchAndStoreItems)
	r.GET("/items", routes.GetAllItems)
	r.GET("/items/:typeID", routes.GetItemByTypeID)

	// New routes
	r.GET("/characters/:id/killmails", routes.GetCharacterKillmails)
	r.GET("/characters/stats", routes.GetAllCharacterStats)

	// New data routes
	r.GET("/characters", routes.GetAllCharacters)
	r.GET("/kills", routes.GetAllKills)

	// Add this line to register the GetKillsByRegion route
	r.GET("/kills/region/:regionID", routes.GetKillsByRegion)

	// Setup Swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Run(":8080")
}
