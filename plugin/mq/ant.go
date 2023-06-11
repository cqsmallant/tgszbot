package mq

import (
	"ant/utils/config"
	"ant/utils/log"

	"github.com/hibiken/asynq"
)

var MClient *asynq.Client

func Start() {
	redis := asynq.RedisClientOpt{
		Addr:     config.RedisDns,
		DB:       config.RedisDb,
		Password: config.RedisPwd,
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
			Concurrency: config.QueueConcurrency,
			// （可选）指定具有不同优先级的多个队列
			Queues: map[string]int{
				"critical": config.QueueLevelCritical,
				"default":  config.QueueLevelDefault,
				"low":      config.QueueLevelLow,
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
