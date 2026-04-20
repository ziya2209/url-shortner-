package config

import (
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

func InitRedisClient() {
	log.Println("idr fir aagya")
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password
		DB:       0,  // use default DB
	})
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	err := rdb.Ping(ctx).Err()
	if err != nil {
		log.Fatal("redis failed to initialize")
	}
	log.Println("succesfully connect with redis")
	RedisClient = rdb

}
