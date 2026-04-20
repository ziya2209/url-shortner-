package main

import (
	"flag"
	"fmt"
	"os"
	"url/config"
	"url/handler"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var dbCon *gorm.DB

func main() {
	// Parse command line flags
	mode := flag.String("mode", "service", "Mode to run: 'service' or 'cron'")
	flag.Parse()

	// Validate mode
	if *mode != "service" && *mode != "cron" {
		fmt.Println("Invalid mode. Use 'service' or 'cron'")
		os.Exit(1)
	}

	// Initialize database connection
	dbCon = config.InitDB()
	config.InitRedisClient()

	// Initialize handler with database connection
	handler.InitHandler(dbCon)

	if *mode == "cron" {
		// Run cron job
		fmt.Println("Running cron job...")
		handler.CornJob()
		fmt.Println("Cron job completed")
		return
	}

	// Run gin server (service mode)
	router := gin.Default()
	router.StaticFile("/", "./frontend/index.html")
	router.POST("/account", handler.RateLimit, handler.CreateAccount)
	router.POST("/login", handler.Login)
	router.POST("/getshorturl", handler.RateLimit, handler.AuthMiddleware, handler.GetShortenUrl)
	router.GET("/:id", handler.RedRedirect)

	fmt.Println("Starting web server on :8081...")
	err := router.Run(":8081")
	if err != nil {
		fmt.Println("your web server falied to run ")
	}
}
