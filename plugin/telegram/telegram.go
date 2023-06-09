package telegram

import (
	"ant/utils/config"
	"ant/utils/log"
	"time"

	tg "gopkg.in/telebot.v3"
)

var Bots *tg.Bot

func BotStart() {
	var err error
	botSetting := tg.Settings{
		Token:  config.TgBotToken,
		Poller: &tg.LongPoller{Timeout: 10 * time.Second},
	}

	if config.ApiProxy != "" {
		botSetting.URL = config.ApiProxy
	}

	Bots, err = tg.NewBot(botSetting)
	if err != nil {
		log.Sugar.Error(err.Error())
		return
	}
	RegisterHandle()
	Bots.Start()

}

func RegisterHandle() {
	// adminOnly := bots.Group()
	// adminOnly.Use(middleware.Whitelist(config.TgManage))

	Bots.Handle(START_CMD, AccountInfo)

	//监听发送文字
	//todo
	Bots.Handle(tg.OnText, OnTextMessageHandle)
	//加群监听，是否机器人，加人
}

func SendToBot(msg string) {
	go func() {
		user := tg.User{
			ID: config.TgManage,
		}
		_, err := Bots.Send(&user, msg, &tg.SendOptions{
			ParseMode: tg.ModeHTML,
		})
		if err != nil {
			log.Sugar.Error(err)
		}
	}()
}

func SendToDice() int {
	diceObj := &tg.Dice{}
	user := tg.User{
		ID: config.TgManage,
	}
	msg, err := diceObj.Send(Bots, &user, &tg.SendOptions{})
	if err != nil {
		log.Sugar.Error(err)
	}
	return msg.Dice.Value

}
