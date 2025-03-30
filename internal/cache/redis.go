package cache

import (
	"context"
	"log"
	"time"

	"github.com/go-redis/redis"
)

var (
	Client *redis.Client
)

func InitRedis(addr, password string) error {
	Client = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0,
	})

	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := Client.Ping().Result()
	if err != nil {
		return err
	}
	log.Println("Connected to Redis")
	return nil
}

func SetFileMetadata(ctx context.Context, fileID string, data interface{}, ttl time.Duration) error {
	return Client.Set("file:"+fileID, data, ttl).Err()
}

func GetFileMetadata(ctx context.Context, fileID string) (string, error) {
	return Client.Get("file:" + fileID).Result()
}

func InvalidateCache(ctx context.Context, key string) error {
	return Client.Del(key).Err()
}
