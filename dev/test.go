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
	// log.Init()
	// dao.MysqlInit()
	// dao.RedisInit()
	// now := time.Now()
	// qs, _ := model.GetQsListByTime(now.Unix())
	// payload, _ := json.Marshal(qs)
	// println(string(payload))
	// ctx := context.Background()
	// dao.Rdb.Set(ctx, constant.CacheQsNow, payload, time.Duration(30)*time.Second)
	// res, _ := dao.Rdb.Get(ctx, constant.CacheQsNow).Result()
	// println(string(res))
	// qs.Status = 3
	// payload2, _ := json.Marshal(qs)
	// println(string(payload2))
	// dao.Rdb.Set(ctx, constant.CacheQsNow, payload2, time.Duration(30)*time.Second)
	// res2, _ := dao.Rdb.Get(ctx, constant.CacheQsNow).Result()
	// println(string(res2))
	// time.Sleep(time.Duration(30) * time.Second)
	// res3, _ := dao.Rdb.Get(ctx, constant.CacheQsNow).Result()
	// println(string(res3))
	// mqtest()
	// diceArr := []int{3, 7, 2}
	// sort.Ints(diceArr)
	// println(diceArr[0])
	// println(diceArr[1])
	// println(diceArr[2])
	println(viper.GetString("api_proxy"))
}

func mqtest() {
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
	qsStartQueue, err := mq.NewQsStartTask(qs)
	if err != nil {
		panic(err)
	}
	mq.MClient.Enqueue(qsStartQueue)

}
