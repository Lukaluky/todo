package db


import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)


var (
	Redis *redis.Client
	Ctx = context.Background()
)

func ConnectRedis() error {
	host := os.Getenv("REDIS_HOST")
	port := os.Getenv("REDIS_PORT")

	if host == "" {
		host = "redis"
	}
	if port == "" {
		port = "6379"
	}

	addr := fmt.Sprintf("%s:%s", host, port)

	Redis = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "", 
		DB:       0,  
	})

	
	ctxTimeout, cancel := context.WithTimeout(Ctx, 3*time.Second)
	defer cancel()

	_, err := Redis.Ping(ctxTimeout).Result()
	if err != nil {
		return fmt.Errorf("failed to connect to Redis: %v", err)
	}

	fmt.Println("Redis connected at", addr)
	return nil
}