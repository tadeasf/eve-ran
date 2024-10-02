package db

import (
	"fmt"
	"log"
	"os"

	"github.com/tadeasf/eve-ran/src/db/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"))
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	fmt.Println("Successfully connected to the database")

	err = InitTables()
	if err != nil {
		log.Fatal("Failed to initialize tables:", err)
	}
}

func InitTables() error {
	return DB.AutoMigrate(
		&models.Character{},
		&models.Kill{},
		&models.Region{},
		&models.System{},
		&models.Constellation{},
		&models.ESIItem{},
	)
}
