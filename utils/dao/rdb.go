package dao

import (
	"ant/utils/config"
	"ant/utils/log"
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gookit/color"
)

var Rdb *redis.Client

func RedisInit() {
	options := redis.Options{
		Addr:        config.RedisDns,                                      // Redis地址
		DB:          config.RedisDb,                                       // Redis库
		PoolSize:    config.RedisPooSize,                                  // Redis连接池大小
		MaxRetries:  config.RedisMaxRetries,                               // 最大重试次数
		IdleTimeout: time.Second * time.Duration(config.RedisIdleTimeout), // 空闲链接超时时间
	}

	if config.RedisPwd != "" {
		options.Password = config.RedisPwd
	}
	Rdb = redis.NewClient(&options)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	pong, err := Rdb.Ping(ctx).Result()
	if err == redis.Nil {
		log.Sugar.Debug("[store_redis] Nil reply returned by Rdb when key does not exist.")
	} else if err != nil {
		color.Red.Printf("[store_redis] redis connRdb err,err=%s", err)
		panic(err)
	} else {
		log.Sugar.Debug("[store_redis] redis connRdb success,suc=%s", pong)
	}
}
