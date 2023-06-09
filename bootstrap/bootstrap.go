package bootstrap

import (
	"ant/command"
	"ant/plugin/mq"
	"ant/utils/dao"
	"ant/utils/log"
)

func Start() {
	defer func() {
		mq.MClient.Close()
	}()
	log.Init()
	dao.MysqlInit()
	dao.RedisInit()
	mq.Start()
	// go telegram.BotStart()
	err := command.Execute()
	if err != nil {
		panic(err)
	}
}
