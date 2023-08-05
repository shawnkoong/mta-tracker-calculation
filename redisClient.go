package main

import (
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

var once sync.Once
var redisClient *redis.Client

func getRedisClient() *redis.Client {
	once.Do(func() {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}

		redisAddress := os.Getenv("REDIS_ADDRESS")
		redisPassword := os.Getenv("REDIS_PASSWORD")
		redisClient = redis.NewClient(&redis.Options{
			Addr:     redisAddress,
			Password: redisPassword,
			DB:       0,
		})
	})
	return redisClient
}
