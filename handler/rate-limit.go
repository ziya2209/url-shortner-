package handler

import (
	"fmt"
	"log"
	"net/http"
	"time"
	"url/config"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func RateLimit(c *gin.Context) {
	ip := c.ClientIP()
	key := "rate_limit:create_account:" + ip

	val, err := config.RedisClient.Get(c, key).Int()
	if err != nil {
		if err == redis.Nil {
			fmt.Println("andr aayi kya request")
			err = config.RedisClient.Set(c, key, 1, time.Minute*15).Err()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "issue comes from server",
				})
				return
			}
			fmt.Println("gyi kya")
			c.Next()
			return

		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "bad request" + err.Error(),
			})
			c.Abort()
			return

		}
	}

	if val < 5 {
		err = config.RedisClient.Incr(c, key).Err()
		if err != nil {
			log.Println("error in incrementing redis value:" + err.Error())
			c.Next()
			return
		}

	} else {
		c.JSON(http.StatusTooManyRequests, gin.H{
			"error": " block for 15 mintues",
		})
		c.Abort()
		return
	}
	c.Next()

}
