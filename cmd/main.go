package main

import (
	"platform-go-challenge/models"

	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func main() {
	InitDB()
	router := gin.Default()

	//routes
	router.POST("/book", CreateBook)
	router.GET("/books", GetBooks)
	router.GET("/book/:id", GetBook)
	router.PUT("/book/:id", UpdateBook)
	router.DELETE("/book/:id", DeleteBook)

	router.Run(":8080")
}

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

	DB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// migrate the schema
	if err := DB.AutoMigrate(&models.Book{}); err != nil {
		log.Fatal("Failed to migrate schema:", err)
	}

	if DB != nil {
		log.Println("DB connection established")
	}
}

func CreateBook(c *gin.Context) {
	var book models.Book

	//bind the request body
	if err := c.ShouldBindJSON(&book); err != nil {
		models.ResponseJSON(c, http.StatusBadRequest, "Invalid input", nil)
		return
	}
	DB.Create(&book)
	models.ResponseJSON(c, http.StatusCreated, "Book created successfully", book)
}

func GetBooks(c *gin.Context) {
	if DB == nil {
		log.Fatal("DB pointer is nil")
	}

	var books []models.Book
	DB.Find(&books)
	models.ResponseJSON(c, http.StatusOK, "Books retrieved successfully", books)
}

func GetBook(c *gin.Context) {
	if DB == nil {
		log.Fatal("DB pointer is nil")
	}

	var book models.Book
	if err := DB.First(&book, c.Param("id")).Error; err != nil {
		models.ResponseJSON(c, http.StatusNotFound, "Book not found", nil)
		return
	}
	models.ResponseJSON(c, http.StatusOK, "Book retrieved successfully", book)
}

func UpdateBook(c *gin.Context) {
	if DB == nil {
		log.Fatal("DB pointer is nil")
	}

	var book models.Book
	if err := DB.First(&book, c.Param("id")).Error; err != nil {
		models.ResponseJSON(c, http.StatusNotFound, "Book not found", nil)
		return
	}

	// bind the request body
	if err := c.ShouldBindJSON(&book); err != nil {
		models.ResponseJSON(c, http.StatusBadRequest, "Invalid input", nil)
		return
	}

	DB.Save(&book)
	models.ResponseJSON(c, http.StatusOK, "Book updated successfully", book)
}

func DeleteBook(c *gin.Context) {
	if DB == nil {
		log.Fatal("DB pointer is nil")
	}

	var book models.Book
	if err := DB.Delete(&book, c.Param("id")).Error; err != nil {
		models.ResponseJSON(c, http.StatusNotFound, "Book not found", nil)
		return
	}
	models.ResponseJSON(c, http.StatusOK, "Book deleted successfully", nil)
}
