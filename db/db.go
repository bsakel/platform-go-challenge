package db

import (
	"log"
	"os"

	"platform-go-challenge/models"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var GormDB *gorm.DB

func InitDB() {

	dsn, isSet := os.LookupEnv("DB_URL")
	if !isSet {
		log.Println("DB_URL environment variable not set, loading from .env file")
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Failed to connect to database:", err)
		}

		dsn = os.Getenv("DB_URL")
	}
	log.Printf("DB_URL value: %s", dsn)

	var err error
	GormDB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// migrate the schema
	if err := GormDB.AutoMigrate(
		&models.Audience{},
		&models.Chart{},
		&models.Insight{},
		&models.UserFavourite{},
	); err != nil {
		log.Fatal("Failed to migrate schema:", err)
	}

	if GormDB != nil {
		log.Println("DB connection established")
	}
}
