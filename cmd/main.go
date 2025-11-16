package main

import (
	"platform-go-challenge/api"

	"github.com/gin-gonic/gin"
)

func main() {
	api.InitDB()
	router := gin.Default()

	//routes
	router.POST("/book", api.CreateBook)
	router.GET("/books", api.GetBooks)
	router.GET("/book/:id", api.GetBook)
	router.PUT("/book/:id", api.UpdateBook)
	router.DELETE("/book/:id", api.DeleteBook)

	router.Run(":8080")
}
