package db

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stakwork/sphinx-tribes/utils"
)

var ctx = context.Background()
var RedisClient *redis.Client
var RedisError error
var expireTime = 6 * time.Hour

func InitRedis() {
	redisURL := os.Getenv("REDIS_URL")

	if redisURL == "" {
		dbInt, _ := utils.ConvertStringToInt(os.Getenv("REDIS_DB"))
		RedisClient = redis.NewClient(&redis.Options{
			Addr:     os.Getenv("REDIS_HOST"),
			Username: os.Getenv("REDIS_USER"),
			Password: os.Getenv("REDIS_PASS"),
			DB:       dbInt,
		})
	} else {
		opt, err := redis.ParseURL(redisURL)
		if err != nil {
			RedisError = err
			fmt.Println("REDIS URL CONNECTION ERROR ===", err)
		}
		RedisClient = redis.NewClient(opt)
	}
	if err := RedisClient.Ping(ctx).Err(); err != nil {
		RedisError = err
		fmt.Println("Could Not Connect To Redis", err)
	}
}

func SetValue(key string, value interface{}) {
	err := RedisClient.Set(ctx, key, value, expireTime).Err()
	if err != nil {
		fmt.Println("REDIS SET ERROR :", err)
	}
}

func GetValue(key string) string {
	val, err := RedisClient.Get(ctx, key).Result()
	if err != nil {
		fmt.Println("REDIS GET ERROR :", err)
	}

	return val
}

func SetMap(key string, values map[string]interface{}) {
	for k, v := range values {
		err := RedisClient.HSet(ctx, key, k, v).Err()
		if err != nil {
			fmt.Println("REDIS SET MAP ERROR :", err)
		}
	}
	RedisClient.Expire(ctx, key, expireTime)
}

func GetMap(key string) map[string]string {
	values := RedisClient.HGetAll(ctx, key).Val()
	return values
}
