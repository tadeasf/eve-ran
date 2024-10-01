package main

import (
	"log"

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
	defer db.DB.Close()

	// Create kills table
	_, err := db.DB.Exec(`
		CREATE TABLE IF NOT EXISTS kills (
			killmail_id BIGINT PRIMARY KEY,
			character_id BIGINT REFERENCES characters(id),
			killmail_time TIMESTAMP,
			solar_system_id INTEGER,
			location_id BIGINT,
			hash TEXT,
			fitted_value NUMERIC,
			dropped_value NUMERIC,
			destroyed_value NUMERIC,
			total_value NUMERIC,
			points INTEGER,
			npc BOOLEAN,
			solo BOOLEAN,
			awox BOOLEAN,
			victim JSONB,
			attackers JSONB
		)
	`)
	if err != nil {
		log.Fatalf("Error creating kills table: %v", err)
	}

	// Start the kill fetcher job
	go jobs.StartKillFetcherJob()

	r := gin.Default()

	// zKillboard routes
	r.POST("/characters", routes.AddCharacter)
	r.DELETE("/characters/:id", routes.RemoveCharacter)
	r.GET("/characters/:id/kills", routes.GetCharacterKills)
	r.GET("/characters/:id/kills/db", routes.GetCharacterKillsFromDB)

	// Setup Swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Run(":8080")
}
