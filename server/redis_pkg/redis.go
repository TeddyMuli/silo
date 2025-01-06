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
	redisAddress := os.Getenv("REDIS_ADDRESS")

	var err error

	options := &redis.Options{
		Addr: redisAddress,
		Username: os.Getenv("REDIS_USER"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
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
