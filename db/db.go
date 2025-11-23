package db

import (
	"log"
	"os"
	"time"

	"platform-go-challenge/models"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var GormDB *gorm.DB

func InitDB() {

	db_url, isSet := os.LookupEnv("DB_URL")
	if !isSet {
		log.Println("DB_URL environment variable not set, loading from .env file")
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Failed to connect to database:", err)
		}

		db_url = os.Getenv("DB_URL")
	}
	log.Printf("DB_URL value: %s", db_url)

	var err error
	GormDB, err = gorm.Open(postgres.Open(db_url), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Configure connection pool for optimal performance
	sqlDB, err := GormDB.DB()
	if err != nil {
		log.Fatal("Failed to get database instance:", err)
	}

	// Set connection pool parameters
	sqlDB.SetMaxIdleConns(10)                  // Maximum idle connections in the pool
	sqlDB.SetMaxOpenConns(100)                 // Maximum open connections to the database
	sqlDB.SetConnMaxLifetime(time.Hour)        // Maximum lifetime of a connection
	sqlDB.SetConnMaxIdleTime(10 * time.Minute) // Maximum idle time before closing

	log.Println("Database connection pool configured: MaxIdle=10, MaxOpen=100, MaxLifetime=1h")

	// migrate the schema
	if err := GormDB.AutoMigrate(
		&models.Audience{},
		&models.Chart{},
		&models.Insight{},
		&models.UserStar{},
	); err != nil {
		log.Fatal("Failed to migrate schema:", err)
	}

	if GormDB != nil {
		log.Println("DB connection established")
	}
}
