package api

import (
	"log"
	"net/http"
	"strconv"

	"platform-go-challenge/api/model"
	"platform-go-challenge/db"
	"platform-go-challenge/models"

	"github.com/gin-gonic/gin"
)

func CreateUserFavourite(c *gin.Context) {
	var userfavourite models.UserFavourite

	//bind the request body
	if err := c.ShouldBindJSON(&userfavourite); err != nil {
		model.ResponseJSON(c, http.StatusBadRequest, "Invalid input", nil)
		return
	}
	db.GormDB.Create(&userfavourite)
	model.ResponseJSON(c, http.StatusCreated, "UserFavourite created successfully", userfavourite)
}

func GetUserFavourites(c *gin.Context) {
	if db.GormDB == nil {
		log.Fatal("DB pointer is nil")
	}

	var userfavourites []models.UserFavourite
	db.GormDB.Find(&userfavourites)
	model.ResponseJSON(c, http.StatusOK, "UserFavourites retrieved successfully", userfavourites)
}

func GetUserFavourite(c *gin.Context) {
	if db.GormDB == nil {
		log.Fatal("DB pointer is nil")
	}

	var userfavourite models.UserFavourite
	if err := db.GormDB.First(&userfavourite, c.Param("id")).Error; err != nil {
		model.ResponseJSON(c, http.StatusNotFound, "UserFavourite not found", nil)
		return
	}
	model.ResponseJSON(c, http.StatusOK, "UserFavourite retrieved successfully", userfavourite)
}

func GetUserFavouritesByUser(c *gin.Context) {
	if db.GormDB == nil {
		log.Fatal("DB pointer is nil")
	}

	userID, err := strconv.Atoi(c.Param("userid"))
	if err != nil {
		model.ResponseJSON(c, http.StatusBadRequest, "Invalid user ID", nil)
		return
	}

	var userfavourites []models.UserFavourite
	if err := db.GormDB.Where("user_id = ?", userID).Find(&userfavourites).Error; err != nil {
		model.ResponseJSON(c, http.StatusNotFound, "UserFavourites not found", nil)
		return
	}
	model.ResponseJSON(c, http.StatusOK, "UserFavourites retrieved successfully", userfavourites)
}

func UpdateUserFavourite(c *gin.Context) {
	if db.GormDB == nil {
		log.Fatal("DB pointer is nil")
	}

	var userfavourite models.UserFavourite
	if err := db.GormDB.First(&userfavourite, c.Param("id")).Error; err != nil {
		model.ResponseJSON(c, http.StatusNotFound, "UserFavourite not found", nil)
		return
	}

	// bind the request body
	if err := c.ShouldBindJSON(&userfavourite); err != nil {
		model.ResponseJSON(c, http.StatusBadRequest, "Invalid input", nil)
		return
	}

	db.GormDB.Save(&userfavourite)
	model.ResponseJSON(c, http.StatusOK, "UserFavourite updated successfully", userfavourite)
}

func DeleteUserFavourite(c *gin.Context) {
	if db.GormDB == nil {
		log.Fatal("DB pointer is nil")
	}

	var userfavourite models.UserFavourite
	if err := db.GormDB.Delete(&userfavourite, c.Param("id")).Error; err != nil {
		model.ResponseJSON(c, http.StatusNotFound, "UserFavourite not found", nil)
		return
	}
	model.ResponseJSON(c, http.StatusOK, "UserFavourite deleted successfully", nil)
}
