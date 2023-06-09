package telegram

import (
	"ant/model"
	"ant/utils/config"
	"ant/utils/log"
	"fmt"
	"regexp"

	"github.com/gookit/goutil/mathutil"
	"github.com/gookit/goutil/strutil"
	tg "gopkg.in/telebot.v3"
)

const (
	jsKfTemp        = "请联系管理员"
	ReplayAddWallet = "请发给我一个合法的钱包地址"

	//第20230318253期\n开奖时间：09:59:30\n投注截至：09:59:00\n➖➖➖➖➖➖➖➖➖➖本期投注\n6杀 - 100 （赔率 1:4）\n➖➖➖➖➖➖➖➖➖➖\n👤玩家：阳光  💰余额：9901
	xzGameStr = "第%s期\n开奖时间：%s\n投注截至：%s\n➖➖➖➖➖➖➖➖➖➖\n本期投注\n%s - %s （赔率 %s）\n➖➖➖➖➖➖➖➖➖➖\n👤玩家：%s  💰余额：%s"
)

// 获取个人信息
func getUser(userData *model.User) (*model.User, error) {
	user, err := model.GetUserInfoByTgId(userData.TgId)
	if err != nil {
		return nil, err
	}
	if user.ID > 0 {
		return user, err
	}
	return model.AddUser(userData)
}

// 更新个人信息
func updateUser(userData *model.User) (*model.User, error) {
	user, err := model.GetUserInfoByTgId(userData.TgId)
	if err != nil {
		return nil, err
	}
	if user.ID > 0 {
		return model.EditUser(userData)
	}
	return model.AddUser(userData)
}

// 下单处理
func xzProcess(ctx tg.Context, user *model.User) error {
	msgObj := ctx.Message()
	senderName := fmt.Sprintf("%s(%s)", user.Nickname, user.Username)
	//大100、小100、单100、双100、大单100、大双100、小单100、小双100、6杀100(表示6点下注100)、对子100、顺子100、豹子100
	txt := msgObj.Text
	xz1Rx := regexp.MustCompile(`^(大|小|单|双|大单|大双|小单|小双|对子|顺子|豹子)(\d+)`)
	match := xz1Rx.FindStringSubmatch(txt)
	if match != nil {
		xzType := match[1]
		xzMoney := match[2]
		xzMoneyF := mathutil.MustFloat(xzMoney)
		if user.Money < xzMoneyF {
			return ctx.Send("账户余额不足，下注失败，请立即充值", &tg.SendOptions{
				ReplyTo: msgObj,
			})
		}
		replyMsg := fmt.Sprintf(xzGameStr, "20230318253", "09:59:30", "09:59:00", xzType, xzMoney, "1:4", senderName, "9901")
		return ctx.Send(replyMsg, &tg.SendOptions{
			ReplyTo: msgObj,
		})
	}
	xz2Rx := regexp.MustCompile(`^(\d)杀(\d+)`)
	match = xz2Rx.FindStringSubmatch(txt)
	if match != nil {
		xzType := match[1]
		xzMoney := match[2]
		xzMoneyF := mathutil.MustFloat(xzMoney)
		if user.Money < xzMoneyF {
			return ctx.Send("账户余额不足，下注失败，请立即充值", &tg.SendOptions{
				ReplyTo: msgObj,
			})
		}
		replyMsg := fmt.Sprintf(xzGameStr, "20230318253", "09:59:30", "09:59:00", xzType+"杀", xzMoney, "1:4", senderName, "9901")
		return ctx.Send(replyMsg, &tg.SendOptions{
			ReplyTo: msgObj,
		})
	}
	return nil
}

// 监听信息
func OnTextMessageHandle(ctx tg.Context) error {
	//群入口
	if ctx.Message().FromGroup() {
		//特定群github配置
		//todo
		sender := ctx.Message().Sender
		userData := &model.User{
			TgId:     strutil.MustString(sender.ID),
			GroupId:  1,
			Username: sender.Username,
			Nickname: fmt.Sprintf("%s%s", sender.FirstName, sender.LastName),
			Status:   "normal",
		}
		user, err := getUser(userData)
		if err != nil {
			log.Sugar.Errorln(err)
		}

		//下注处理
		xzProcess(ctx, user)

	}

	if ctx.Message().Private() {
		//个人聊天
		return ctx.Send(jsKfTemp, &tg.SendOptions{
			ParseMode: tg.ModeHTML,
		})
	}
	return nil
}

// ========================================================
// ======================个人处理==========================
// ========================================================
// 个人页面->个人页面按钮
func fnPrivteLnKeyBoard() [][]tg.InlineButton {
	accountInfoBtn := tg.InlineButton{
		Text:   "👤账户信息",
		Unique: "accountInfo",
	}
	Bots.Handle(&accountInfoBtn, AccountInfo)

	gameInfoBtn := tg.InlineButton{
		Text:   "🎮玩法介绍",
		Unique: "gameInfo",
	}
	Bots.Handle(&gameInfoBtn, GameInfo)

	rechargeBtn := tg.InlineButton{
		Text:   "💰充值入金",
		Unique: "recharge",
	}
	Bots.Handle(&rechargeBtn, fnRecharge)

	withdrawalBtn := tg.InlineButton{
		Text:   "🍮提现出金",
		Unique: "withdrawal",
	}
	Bots.Handle(&withdrawalBtn, fnWithdrawal)

	synAccountBtn := tg.InlineButton{
		Text:   "⚙️账户设置",
		Unique: "synAccount",
	}
	Bots.Handle(&synAccountBtn, fnSynAccount)

	btns := [][]tg.InlineButton{{accountInfoBtn, gameInfoBtn}, {rechargeBtn, withdrawalBtn}, {synAccountBtn}}
	return btns
}

// 个人页面->👤账户信息
func AccountInfo(ctx tg.Context) error {
	chat := ctx.Chat()
	if chat.Type != "private" {
		return nil
	}
	userData := &model.User{
		TgId:     strutil.MustString(chat.ID),
		GroupId:  1,
		Username: chat.Username,
		Nickname: fmt.Sprintf("%s%s", chat.FirstName, chat.LastName),
		Status:   "normal",
	}
	configList, err := model.ConfigList()
	if err != nil {
		return ctx.Send(jsKfTemp, &tg.SendOptions{
			ParseMode: tg.ModeHTML,
		})
	}
	wzName := configList[0].Value
	wzUrl := configList[20].Value

	user, err := getUser(userData)
	if err != nil {
		return ctx.Send(jsKfTemp, &tg.SendOptions{
			ParseMode: tg.ModeHTML,
		})
	}
	msgTemp := "👤账户信息\n\n【<a href='%s'>%s</a>】Telegram 官方骰子，具体玩法看置顶\n\n"
	msgTemp += "【账户ID】：<code>%s</code>\n【账户昵称】：<b>%s</b>\n【账户余额】：<span class='tg-spoiler'>%f</span>\n"
	return ctx.EditOrSend(fmt.Sprintf(msgTemp, wzUrl, wzName, user.TgId, user.Nickname, user.Money), &tg.SendOptions{
		ParseMode:   tg.ModeHTML,
		ReplyMarkup: &tg.ReplyMarkup{InlineKeyboard: fnPrivteLnKeyBoard()},
	})
}

// 个人页面->🎮玩法介绍
func GameInfo(ctx tg.Context) error {
	configList, err := model.ConfigList()
	if err != nil {
		return ctx.Send(jsKfTemp, &tg.SendOptions{
			ParseMode: tg.ModeHTML,
		})
	}
	wzName := configList[0].Value
	wzUrl := configList[20].Value

	dxds_rate := configList[21].Value
	dddsxdxs_rate := configList[22].Value

	dz_rate := configList[23].Value
	sz_rate := configList[24].Value
	bz_rate := configList[25].Value

	d3_rate := configList[26].Value
	d4_rate := configList[27].Value
	d5_rate := configList[28].Value
	d6_rate := configList[29].Value
	d7_rate := configList[30].Value
	d8_rate := configList[31].Value
	d9_rate := configList[32].Value
	d10_rate := configList[33].Value
	d11_rate := configList[34].Value
	d12_rate := configList[35].Value
	d13_rate := configList[36].Value
	d14_rate := configList[37].Value
	d15_rate := configList[38].Value
	d16_rate := configList[39].Value
	d17_rate := configList[40].Value
	d18_rate := configList[41].Value

	msgTemp := "🎮玩法介绍\n\n【<a href='%s'>%s</a>】Telegram 官方骰子，具体玩法看置顶\n\n"
	msgTemp += "【<a href='%s'>%s</a>】规则\n"

	msgTemp += "大：投掷点数大于等于11（赔率1:%s）\n小：投掷点数小于等于10（赔率1:%s）\n"
	msgTemp += "单：投掷点数为单（赔率1:%s）\n双：投掷点数为双（赔率1:%s）\n"

	msgTemp += "大单：投掷点数大于等于11且点数为单（赔率1:%s）\n大双：投掷点数大于等于11且点数为双（赔率1:%s）\n"
	msgTemp += "小单：投掷点数小于等于10且点数为单（赔率1:%s）\n小双：投掷点数小于等于10且点数为双（赔率1:%s）\n"

	msgTemp += "对子：投掷点数中有两个相同点数（赔率1:%s）\n顺子：投掷点数为顺数（赔率1:%s）\n豹子： 三颗骰子点数相同（赔率1:%s）\n"
	msgTemp += "点杀：投掷点数等于下注点数，赔率如下：\n"
	msgTemp += "3点（赔率1:%s）\n4点（赔率1:%s）\n5点（赔率1:%s）\n6点（赔率1:%s）\n7点（赔率1:%s）\n8点（赔率1:%s）\n9点（赔率1:%s）\n"

	msgTemp += "10点（赔率1:%s）\n11点（赔率1:%s）\n12点（赔率1:%s）\n13点（赔率1:%s）\n14点（赔率1:%s）\n15点（赔率1:%s）\n16点（赔率1:%s）\n"
	msgTemp += "17点（赔率1:%s）\n18点（赔率1:%s）\n\n点击进入【<a href='%s'>%s</a>】\n"
	msgTemp += "发送投注内容：下注类型+金额，例如：\n"
	msgTemp += "<code>大100</code>、<code>小100</code>、<code>单100</code>、<code>双100</code>、<code>大单100</code>、<code>大双100</code>、"
	msgTemp += "<code>小单100</code>、<code>小双100</code>、<code>6杀100(表示6点下注100)</code>、<code>对子100</code>、<code>顺子100</code>、<code>豹子100</code>"

	return ctx.EditOrSend(fmt.Sprintf(msgTemp, wzUrl, wzName, wzUrl, wzName,
		dxds_rate, dxds_rate, dxds_rate, dxds_rate,
		dddsxdxs_rate, dddsxdxs_rate, dddsxdxs_rate, dddsxdxs_rate,
		dz_rate, sz_rate, bz_rate,
		d3_rate, d4_rate, d5_rate, d6_rate, d7_rate, d8_rate, d9_rate,
		d10_rate, d11_rate, d12_rate, d13_rate, d14_rate, d15_rate, d16_rate,
		d17_rate, d18_rate, wzUrl, wzName), &tg.SendOptions{
		ParseMode:   tg.ModeHTML,
		ReplyMarkup: &tg.ReplyMarkup{InlineKeyboard: fnPrivteLnKeyBoard()},
	})
}

// 个人页面->⚙️账户设置
func fnSynAccount(ctx tg.Context) error {
	chat := ctx.Chat()
	if chat.Type != "private" {
		return nil
	}

	userData := &model.User{
		TgId:     strutil.MustString(chat.ID),
		GroupId:  1,
		Username: chat.Username,
		Nickname: fmt.Sprintf("%s%s", chat.FirstName, chat.LastName),
	}
	user, err := updateUser(userData)
	if err != nil {
		return ctx.Send(jsKfTemp, &tg.SendOptions{
			ParseMode: tg.ModeHTML,
		})
	}
	msgTemp := "⚙️账户设置\n\n"
	msgTemp += "【账户ID】：<code>%s</code>\n【账户昵称】：<b>%s</b>\n\n您的账户信息已经自动同步完成，无需进行其他设置！"
	return ctx.EditOrSend(fmt.Sprintf(msgTemp, user.TgId, user.Nickname), &tg.SendOptions{
		ParseMode:   tg.ModeHTML,
		ReplyMarkup: &tg.ReplyMarkup{InlineKeyboard: fnPrivteLnKeyBoard()},
	})
}

// ======================充值页面==========================

// 充值页面->充值页面按钮
func fnRechargeLnKeyBoard() [][]tg.InlineButton {
	configList, _ := model.ConfigList()
	kfUrl := configList[19].Value
	rechargeBtn := tg.InlineButton{
		Text:   "💰充值",
		Unique: "rechargeUrl",
		URL:    kfUrl,
	}

	rechargeListBtn := tg.InlineButton{
		Text:   "🗓记录",
		Unique: "rechargeList",
	}
	Bots.Handle(&rechargeListBtn, fnRechargeList)

	backBtn := tg.InlineButton{
		Text:   "⤴️返回",
		Unique: "rechargeBack",
	}
	Bots.Handle(&backBtn, AccountInfo)

	btns := [][]tg.InlineButton{{rechargeBtn, rechargeListBtn}, {backBtn}}
	return btns
}

// 充值页面->💰充值入金
func fnRecharge(ctx tg.Context) error {
	msgTemp := "💰充值入金\n\n收款钱包地址为：<code>%s</code>\n（点击复制）\n充值完成后请点击下方充值按钮联系客服确认充值"
	return ctx.EditOrSend(fmt.Sprintf(msgTemp, config.TrcToken), &tg.SendOptions{
		ParseMode:   tg.ModeHTML,
		ReplyMarkup: &tg.ReplyMarkup{InlineKeyboard: fnRechargeLnKeyBoard()},
	})
}

// 充值页面->🗓记录
func fnRechargeList(ctx tg.Context) error {
	msgTemp := "近10次充值记录\n%s"
	list := "无充值记录"
	return ctx.Reply(fmt.Sprintf(msgTemp, list), &tg.SendOptions{
		ParseMode: tg.ModeHTML,
	})
}

// ======================提现页面==========================
// 提现页面->提现页面按钮
func fnWithdrawalLnKeyBoard() [][]tg.InlineButton {
	configList, _ := model.ConfigList()
	kfUrl := configList[19].Value
	withdrawalBtn := tg.InlineButton{
		Text:   "🍮提现",
		Unique: "withdrawalUrl",
		URL:    kfUrl,
	}

	withdrawalListBtn := tg.InlineButton{
		Text:   "🗓记录",
		Unique: "withdrawalList",
	}
	Bots.Handle(&withdrawalListBtn, fnWithdrawalList)

	backBtn := tg.InlineButton{
		Text:   "⤴️返回",
		Unique: "withdrawalBack",
	}
	Bots.Handle(&backBtn, AccountInfo)

	btns := [][]tg.InlineButton{{withdrawalBtn, withdrawalListBtn}, {backBtn}}
	return btns
}

// 提现页面->🍮提现出金
func fnWithdrawal(ctx tg.Context) error {
	msgTemp := "🍮提现出金\n\n如需提现请点击下方提现按钮联系客服，并说明提现金额和提供收款钱包地址即可"
	return ctx.EditOrSend(msgTemp, &tg.SendOptions{
		ParseMode:   tg.ModeHTML,
		ReplyMarkup: &tg.ReplyMarkup{InlineKeyboard: fnWithdrawalLnKeyBoard()},
	})
}

// 提现页面->🗓记录
func fnWithdrawalList(ctx tg.Context) error {
	msgTemp := "近10次提现记录\n%s"
	list := "无提现记录"
	return ctx.Reply(fmt.Sprintf(msgTemp, list), &tg.SendOptions{
		ParseMode: tg.ModeHTML,
	})
}
