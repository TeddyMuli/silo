package redis_pkg

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

// Initializa a redis connection
func InitRedis() {
	env := os.Getenv("ENV")

	var redisAddress string
	if env == "production" {
		redisAddress = "redis://redis:6379"
	} else {
		redisAddress = "127.0.0.1:6379"
	}

	var options *redis.Options
	var err error

	if env == "production" {
		options, err = redis.ParseURL(redisAddress)
		if err != nil {
			fmt.Println("Error parsing Redis URL:", err)
			os.Exit(1)
		}
	} else {
		options = &redis.Options{
			Addr: redisAddress,
		}
	}

	RedisClient = redis.NewClient(options)
	
	// Test the connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = RedisClient.Ping(ctx).Result()
	if err != nil {
		fmt.Println("Could not connect to Redis:", err)
	} else {
		fmt.Println("Connected to Redis!")
	}
}
