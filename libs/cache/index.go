package cache

import (
	"context"
	"time"

	"github.com/charmbracelet/log"
	"github.com/redis/go-redis/v9"
)

func GetClient(inAddr string, inPassword string, inDb int) *redis.Client {
	log.Info("Init redis ...")
	// Get Redis connection info from environment variables
	redisClient := redis.NewClient(&redis.Options{
		Addr:         inAddr,
		Password:     inPassword,
		DB:           inDb,
		PoolSize:     10,              // Default connection pool size
		MinIdleConns: 5,               // Minimum number of idle connections
		MaxRetries:   3,               // Maximum number of retries
		DialTimeout:  5 * time.Second, // Timeout for establishing new connections
		ReadTimeout:  3 * time.Second, // Timeout for socket reads
		WriteTimeout: 3 * time.Second, // Timeout for socket writes
		PoolTimeout:  4 * time.Second, // Timeout for getting connection from pool
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := redisClient.Ping(ctx).Err(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	log.Info("Connected to Redis success")
	return redisClient
}
