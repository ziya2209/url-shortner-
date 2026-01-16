package handler

import (
	"net/http"
	"strings"
	"time"
	"url/model"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var dbCon *gorm.DB

// InitHandler initializes the database connection for handlers
func InitHandler(db *gorm.DB) {
	dbCon = db
}

func CreateAccount(c *gin.Context) {
	var req model.AccountCreationData
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var user *model.User
	err = dbCon.Where("email = ?", req.Email).First(&user).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return

	}
	if user.Email == req.Email {
		c.JSON(http.StatusConflict, gin.H{
			"error": "this email already exists",
		})
		return
	}

	req.Email = strings.TrimSpace(req.Email)
	req.Name = strings.TrimSpace(req.Name)
	req.Password = strings.TrimSpace(req.Password)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	user = &model.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
	}
	err = dbCon.Create(user).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}
	c.JSON(200, req)

}

func Login(c *gin.Context) {
	var req model.AccountLogingData
	err := c.ShouldBindBodyWithJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	var user *model.User

	err = dbCon.Where("email = ?", req.Email).First(&user).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return

	}
	if user.Email != req.Email {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "email not found",
		})
		return

	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "invalid password",
		})
		return
	}

	secret := "ziya"

	claims := jwt.MapClaims{
		"name":  user.Name,
		"email": user.Email,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to sigh the token",
		})
		return

	}
	c.JSON(200, gin.H{
		"msg":   "user login succesfully",
		"token": tokenString,
	})


}
