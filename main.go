package main

import (
	"fmt"
	"url/config"
	"url/handler"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var dbCon *gorm.DB

func main() {
	// Initialize database connection
	dbCon = config.InitDB()

	// Initialize handler with database connection
	handler.InitHandler(dbCon)

	router := gin.Default()
	router.POST("/account", handler.CreateAccount)
	router.POST("/login", handler.Login)
	err := router.Run(":8081")
	if err != nil {
		fmt.Println("your web server falied to run ")
	}

}
