package util

import (
	"fmt"
	"github.com/go-redis/redis/v7"
	"log"
	"time"
)

var rc *redis.Client

func init() {
	NewClient()
}

func NewClient() {
	rc = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:6379", IP),
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

func SetNX(key string, value interface{}, expiration time.Duration) (bool, error) {
	ok, err := rc.SetNX(key, value, expiration).Result()
	if err != nil {
		return false, err
	}
	if ok {
		// Key was set successfully because it did not exist
		return true, err
	} else {
		// Key was not set because it already exists
		return false, err
	}
}
