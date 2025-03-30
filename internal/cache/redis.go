package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	RedisClient *redis.Client
	ctx         = context.Background()
)

func InitRedis(addr string) error {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	_, err := RedisClient.Ping(ctx).Result()
	return err
}

func SetFileMetadata(fileID string, file interface{}, expiration time.Duration) error {
	return RedisClient.Set(ctx, "file:"+fileID, file, expiration).Err()
}

func GetFileMetadata(fileID string) (string, error) {
	return RedisClient.Get(ctx, "file:"+fileID).Result()
}

func InvalidateFileCache(fileID string) error {
	return RedisClient.Del(ctx, "file:"+fileID).Err()
}
