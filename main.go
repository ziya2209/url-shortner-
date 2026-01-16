package main

import (
	"url/handler"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)


var dbCon *gorm.DB

func main(){
	router := gin.Default()
	router.POST("/account",handler.CreateAccount)
	router.POST("/login",handler.Login)

}