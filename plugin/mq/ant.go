package mq

import (
	"ant/utils/log"
	"fmt"

	"github.com/hibiken/asynq"
	"github.com/spf13/viper"
)

var MClient *asynq.Client

func Start() {
	redis := asynq.RedisClientOpt{
		Addr: fmt.Sprintf(
			"%s:%s",
			viper.GetString("redis_host"),
			viper.GetString("redis_port")),
		DB:       viper.GetInt("redis_db"),
		Password: viper.GetString("redis_passwd"),
	}
	InitClient(redis)
	go initServe(redis)
}

func InitClient(redis asynq.RedisClientOpt) {
	MClient = asynq.NewClient(redis)
}

func initServe(redis asynq.RedisClientOpt) {
	srv := asynq.NewServer(redis,
		asynq.Config{
			// 每个进程并发执行的worker数量
			Concurrency: viper.GetInt("queue_concurrency"),
			// （可选）指定具有不同优先级的多个队列
			Queues: map[string]int{
				"critical": viper.GetInt("queue_level_critical"),
				"default":  viper.GetInt("queue_level_default"),
				"low":      viper.GetInt("queue_level_low"),
			},
			Logger: log.Sugar,
		})
	mux := asynq.NewServeMux()
	mux.HandleFunc(QsStart, QsStartHandle)
	mux.HandleFunc(QsLocking, QsLockingHandle)
	mux.HandleFunc(QsLocked, QsLockedHandle)
	mux.HandleFunc(QsClosed, QsClosedHandle)
	if err := srv.Run(mux); err != nil {
		log.Sugar.Fatal("[queue] could not run server: %v", err)
	}
}
