package main

import (
	"ant/model"
	"ant/plugin/mq"
	"ant/utils/dao"
	"ant/utils/log"
	"fmt"
	"time"

	"github.com/hibiken/asynq"
	"github.com/spf13/viper"
)

func main() {
	log.Init()
	dao.MysqlInit()
	now := time.Now()
	qs, _ := model.GetQsListByTime(now.Unix())
	println(qs.ID)
	redis := asynq.RedisClientOpt{
		Addr: fmt.Sprintf(
			"%s:%s",
			viper.GetString("redis_host"),
			viper.GetString("redis_port")),
		DB:       viper.GetInt("redis_db"),
		Password: viper.GetString("redis_passwd"),
	}
	mq.InitClient(redis)

	// 超时过期消息队列
	qsStartQueue, _ := mq.NewQsStartTask(qs)
	mq.MClient.Enqueue(qsStartQueue)
	//mq.MClient.Enqueue(orderExpirationQueue, asynq.ProcessIn(config.GetOrderExpirationTimeDuration()))

	// ExpirationTime := carbon.Now().AddMinutes(config.GetOrderExpirationTime()).Timestamp()

}
