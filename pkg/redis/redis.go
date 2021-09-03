package redis

import (
	"fork_go_im/pkg/config"
	"log"
	"time"

	"github.com/go-redis/redis"
)

var RedisDB *redis.Client

func InitClient() (err error) {
	RedisDB = redis.NewClient(&redis.Options{
		Network:      "tcp",
		Addr:         config.GetString("cache.redis.addr") + ":" + config.GetString("cache.redis.port"),
		Password:     config.GetString("cache.redis.password"),
		DB:           config.GetInt("cache.redis.db", 0),
		PoolSize:     15,
		MinIdleConns: 10,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		PoolTimeout:  5 * time.Second,
	})

	_, err = RedisDB.Ping().Result()
	if err != nil {
		log.Println(err)
	}
	return nil
}
