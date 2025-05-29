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
		Password: "Bwin@0913",
		DB:       0,
	})
	r, err := rc.Ping().Result()
	if err != nil {
		log.Fatalf("Redis连接失败: %v", err)
	}
	log.Println("Pong:", r)
}

func SAdd(k, m string) int64 {
	r, err := rc.SAdd(k, m).Result()
	if err != nil {
		log.Printf("SAdd Error: %v\n", err)
		return -1
	}
	return r
}

func SPop(key string) string {
	r, err := rc.SPop(key).Result()
	if err != nil {
		log.Printf("SPop Error: %v\n", err)
		return ""
	}
	return r
}

func SIsMember(key, m string) (bool, error) {
	r, err := rc.SIsMember(key, m).Result()
	if err != nil {
		return false, err
	}
	return r, nil
}

func SetNX(key string, value interface{}, expiration time.Duration) bool {
	ok, err := rc.SetNX(key, value, expiration).Result()
	if err != nil {
		log.Printf("SetNX Error: %v\n", err)
		return false
	}
	if ok {
		// Key was set successfully because it did not exist
		return true
	} else {
		// Key was not set because it already exists
		return false
	}
}

// RPush 向 List 中添加元素
func RPush(key string, value interface{}) (int64, error) {
	result, err := rc.RPush(key, value).Result()
	if err != nil {
		return 0, err
	}
	return result, err
}

// LPop 从 List 中移除第一个元素
func LPop(key string) (string, error) {
	result, err := rc.LPop(key).Result()
	if err != nil {
		return "", err
	}
	return result, err
}

// Del 删除 key
func Del(key string) (int64, error) {
	result, err := rc.Del(key).Result()
	if err != nil {
		return 0, err
	}
	return result, nil
}

// LLen 获取列表的长度
func LLen(key string) (int64, error) {
	result, err := rc.LLen(key).Result()
	if err != nil {
		return 0, err
	}
	return result, nil
}
