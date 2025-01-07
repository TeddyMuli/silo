package otp

import (
	"context"
	"fmt"
	"time"

	"server/redis_pkg"

	"github.com/redis/go-redis/v9"
)

func StoreOTP(email string, otp string) error {
	// Store the OTP with an expiration time (e.g., 5 minutes)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err := redis_pkg.RedisClient.Set(ctx, email+"_otp", otp, time.Minute*5).Err()
	if err != nil {
		return fmt.Errorf("failed to store OTP: %w", err)
	}

	return nil
}

func GetStoredOTP(email string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	otp, err := redis_pkg.RedisClient.Get(ctx, email+"_otp").Result()
	if err == redis.Nil {
		return "", fmt.Errorf("OTP not found or expired")
	} else if err != nil {
		return "", fmt.Errorf("failed to retrieve OTP: %w", err)
	}

	return otp, nil
}
