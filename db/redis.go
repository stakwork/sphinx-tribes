package db

import (
	"context"
	"fmt"
	"os"

	"github.com/redis/go-redis/v9"
	"github.com/stakwork/sphinx-tribes/utils"
)

var ctx = context.Background()
var RedisClient *redis.Client

func InitRedis() {
	redisURL := os.Getenv("REDIS_URL")
	fmt.Println("redis url :", redisURL)

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
			fmt.Println("REDIS URL CONNECTION ERROR ===", err)
		}

		RedisClient = redis.NewClient(opt)
	}
}

func SetValue(key string, value string) {
	err := RedisClient.Set(ctx, key, value, 0).Err()
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

func SetMap(key string, values map[string]string) {
	for k, v := range values {
		err := RedisClient.HSet(ctx, key, k, v).Err()
		if err != nil {
			fmt.Println("REDIS SET MAP ERROR :", err)
		}
	}
}

func GetMap(key string) map[string]string {
	values := RedisClient.HGetAll(ctx, key).Val()
	return values
}
