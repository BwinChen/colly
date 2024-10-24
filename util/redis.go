package util

import (
	"github.com/go-redis/redis/v7"
	"log"
)

var rc *redis.Client

func init() {
	NewClient()
}

func NewClient() {
	rc = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	r, err := rc.Ping().Result()
	if err != nil {
		log.Fatalf("Redis连接失败: %v", err)
	}
	log.Println("Pong:", r)
}

func SAdd(id string) (int64, error) {
	r, err := rc.SAdd("ids", id).Result()
	if err != nil {
		return 0, err
	}
	return r, nil
}

func SIsMember(m string) (bool, error) {
	r, err := rc.SIsMember("ids", m).Result()
	if err != nil {
		return false, err
	}
	return r, nil
}

func SRem(m string) (int64, error) {
	r, err := rc.SRem("ids", m).Result()
	if err != nil {
		return 0, err
	}
	return r, err
}
