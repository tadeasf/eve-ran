package db

import (
	"fmt"
	"log"

	"github.com/tadeasf/eve-ran/src/db/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	dsn := "host=localhost port=5435 user=eve password=eve dbname=eve sslmode=disable"
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	fmt.Println("Successfully connected to the database")

	// Auto Migrate the schema
	migrateSchema()
}

func migrateSchema() {
	DB.AutoMigrate(
		&models.Character{},
		&models.Kill{},
		&models.Region{},
		&models.System{},
		&models.Constellation{},
		&models.ESIItem{},
	)
}
