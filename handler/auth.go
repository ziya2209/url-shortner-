package handler

import (
	"fmt"
	"net/http"
	"strings"
	"time"
	"url/helper"
	"url/model"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

func AuthMiddleware(c *gin.Context) {

	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "Authorization header is missing",
		})
		return
	}
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "Authorization header format must be Bearer {token}",
		})
		return
	}
	tokenString := parts[1]

	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok || t.Method.Alg() != "HS256" {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(secret), nil
	})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "invalid token",
			"msg":   err.Error(),
		})
		return
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid claims",
		})
		return
	}
	if exp, ok := claims["exp"].(float64); ok {
		if int64(exp) < time.Now().Unix() {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "token expired",
			})
			return
		}
	}

	email, ok := claims["email"].(string)
	if !ok || email == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "email is not present in the token",
		})
		return

	}
	var user model.User
	err = dbCon.Where("email = ?", email).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "email not found",
			})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "database error" + err.Error(),
		})
		return

	}
	c.Set("user", user) // store claims and context

	c.Next()

}

func GetShortenUrl(c *gin.Context) {
	var req model.LongUrl

	val, ok := c.Get("user")
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "context not updated",
		})
		return
	}
	user := val.(model.User)

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	shortId := helper.GenerateShortUrl(fmt.Sprintf("%d", user.ID), req.LongURL)

	url := model.URL{
		LongURL:    req.LongURL,
		ShortURLID: shortId,
		UserID:     int64(user.ID),
	}
	err = dbCon.Create(&url).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create url," + err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"id": url.ShortURLID,
	})

}

func RedRedirect(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "short_id cannot be empty",
		})
		return
	}
	var url model.URL
	err := dbCon.Where("short_url_id = ?", id).First(&url).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "something went wrong or short id does not exist",
		})
		return
	}
	go func() {
		var ClickCheck model.URLClicks
		err = dbCon.Where("url_id = ?", url.ID).First(&ClickCheck).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "something went wrong or short id does not exist",
			})
			return

		}
		if err != nil {
			ClickCheck.ClickedAt = time.Now()
			ClickCheck.URLId = url.ID
			ClickCheck.IPAddress = c.RemoteIP()

			err = dbCon.Create(&ClickCheck).Error
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
				return

			}
		} else {
			//ClickCheck.ClickedAt = time.Now()
			//ClickCheck.IPAddress = c.RemoteIP()
			err = dbCon.Model(&model.URLClicks{}).Where("url_id = ?", url.ID).Updates(map[string]any{
				"clicked_at": time.Now(),
				"ip_address": c.RemoteIP(),
			}).Error
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
				return

			}

		}
	}()

	c.Redirect(301, url.LongURL)

}
