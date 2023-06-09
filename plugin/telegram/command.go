package telegram

import tg "gopkg.in/telebot.v3"

const (
	START_CMD = "/start"
)

var Cmds = []tg.Command{
	{
		Text:        START_CMD,
		Description: "开始",
	},
}
