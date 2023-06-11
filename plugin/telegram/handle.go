package telegram

import (
	"ant/model"
	"ant/utils/config"
	"ant/utils/constant"
	"ant/utils/dao"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gookit/goutil/mathutil"
	"github.com/gookit/goutil/strutil"
	tg "gopkg.in/telebot.v3"
)

const (
	gameLockedTemp  = "已封盘，等待下一期"
	moneyBZTemp     = "账户余额不足，下注失败，请立即充值"
	billRemarkTemp  = "%s[第%s期下注]%s - %.2f （赔率 1:%.2f）"
	jsKfTemp        = "请联系管理员"
	ReplayAddWallet = "请发给我一个合法的钱包地址"
	//📣小薛 的充值10000已成功到账！
	//第20230318253期\n开奖时间：09:59:30\n投注截至：09:59:00\n➖➖➖➖➖➖➖➖➖➖本期投注\n6杀 - 100 （赔率 1:4）\n➖➖➖➖➖➖➖➖➖➖\n👤玩家：阳光  💰余额：9901
	xzGameStr  = "<b>第%s期\n开奖时间：%s\n投注截至：%s\n➖➖➖➖➖➖➖➖➖➖\n本期投注</b>\n%s➖➖➖➖➖➖➖➖➖➖\n👤玩家：%s  💰余额：<code>%.2f</code>\n"
	xzGameStr2 = "<code>%s - %.2f</code> （赔率 1:%.2f）\n"
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

// 监听信息
func xzInfoBytxTypeStr(xzType string, configList []model.Config) (int, float64) {

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
	var stake int = 0
	var rate float64 = 0
	switch xzType {
	case "大":
		rate = mathutil.MustFloat(dxds_rate)
		stake = 1
	case "小":
		rate = mathutil.MustFloat(dxds_rate)
		stake = 2

	case "3":
		rate = mathutil.MustFloat(d3_rate)
		stake = 3
	case "4":
		rate = mathutil.MustFloat(d4_rate)
		stake = 4
	case "5":
		rate = mathutil.MustFloat(d5_rate)
		stake = 5
	case "6":
		rate = mathutil.MustFloat(d6_rate)
		stake = 6
	case "7":
		rate = mathutil.MustFloat(d7_rate)
		stake = 7
	case "8":
		rate = mathutil.MustFloat(d8_rate)
		stake = 8
	case "9":
		rate = mathutil.MustFloat(d9_rate)
		stake = 9
	case "10":
		rate = mathutil.MustFloat(d10_rate)
		stake = 10
	case "11":
		rate = mathutil.MustFloat(d11_rate)
		stake = 11
	case "12":
		rate = mathutil.MustFloat(d12_rate)
		stake = 12
	case "13":
		rate = mathutil.MustFloat(d13_rate)
		stake = 13
	case "14":
		rate = mathutil.MustFloat(d14_rate)
		stake = 14
	case "15":
		rate = mathutil.MustFloat(d15_rate)
		stake = 15
	case "16":
		rate = mathutil.MustFloat(d16_rate)
		stake = 16
	case "17":
		rate = mathutil.MustFloat(d17_rate)
		stake = 17
	case "18":
		rate = mathutil.MustFloat(d18_rate)
		stake = 18
	case "单":
		rate = mathutil.MustFloat(dxds_rate)
		stake = 19
	case "双":
		rate = mathutil.MustFloat(dxds_rate)
		stake = 20
	case "大单":
		rate = mathutil.MustFloat(dddsxdxs_rate)
		stake = 21
	case "大双":
		rate = mathutil.MustFloat(dddsxdxs_rate)
		stake = 22
	case "小单":
		rate = mathutil.MustFloat(dddsxdxs_rate)
		stake = 23
	case "小双":
		rate = mathutil.MustFloat(dddsxdxs_rate)
		stake = 24
	case "对子":
		rate = mathutil.MustFloat(dz_rate)
		stake = 25
	case "顺子":
		rate = mathutil.MustFloat(sz_rate)
		stake = 26
	case "豹子":
		rate = mathutil.MustFloat(bz_rate)
		stake = 27

	}
	return stake, rate
}

func xzInfoBytxTypeInt(xzType int) string {
	var str string = ""
	if xzType == 1 {
		str = "大"
	} else if xzType == 2 {
		str = "小"
	} else if xzType >= 3 && xzType <= 18 {
		str = fmt.Sprintf("%d杀", xzType)
	} else if xzType == 19 {
		str = "单"
	} else if xzType == 20 {
		str = "双"
	} else if xzType == 21 {
		str = "大单"
	} else if xzType == 22 {
		str = "大双"
	} else if xzType == 23 {
		str = "小单"
	} else if xzType == 24 {
		str = "小双"
	} else if xzType == 25 {
		str = "对子"
	} else if xzType == 26 {
		str = "顺子"
	} else if xzType == 27 {
		str = "豹子"
	}
	return str
}

// 检查封盘
func fnQsLocked(ctx tg.Context) (model.Qs, error) {
	var qs model.Qs
	msgObj := ctx.Message()
	//是否可以下注
	qsPayload, err := dao.Rdb.Get(context.Background(), constant.CacheQsNow).Result()
	if err == redis.Nil || err != nil {
		ctx.Send(gameLockedTemp, &tg.SendOptions{
			ReplyTo: msgObj,
		})
		return qs, err
	}
	err = json.Unmarshal([]byte(qsPayload), &qs)
	if err != nil || qs.Status != 1 {
		ctx.Send(gameLockedTemp, &tg.SendOptions{
			ReplyTo: msgObj,
		})
		return qs, errors.New("已封盘")
	}
	return qs, nil
}
func OnTextMessageHandle(ctx tg.Context) error {

	//群入口
	if ctx.Message().FromGroup() {
		msgObj := ctx.Message()
		txt := msgObj.Text
		var err error
		//特定群配置
		//todo
		authRule, err := model.GetAuthRuleById(40)
		if err != nil {
			return err
		}
		if (-1 * authRule.CreateTime) != ctx.Chat().ID {
			return ctx.Send(jsKfTemp+"<a href='https://t.me/sunnant'>@技术</a>", &tg.SendOptions{
				ReplyTo:   msgObj,
				ParseMode: tg.ModeHTML,
			})
		}
		kfUrlConfig, _ := model.GetConfigByName("kf_url")
		kfNameConfig, _ := model.GetConfigByName("kf_name")

		//充值
		tzRx := regexp.MustCompile(`^(充值)(\d*)`)
		match := tzRx.FindStringSubmatch(txt)
		if match != nil {
			return ctx.Send(fmt.Sprintf("%s<a href='%s'>@%s</a>", jsKfTemp, kfUrlConfig.Value, kfNameConfig.Value), &tg.SendOptions{
				ReplyTo:   msgObj,
				ParseMode: tg.ModeHTML,
			})
		}
		//提现
		txRx := regexp.MustCompile(`^(提现)(\d*)`)
		match = txRx.FindStringSubmatch(txt)
		if match != nil {
			return ctx.Send(fmt.Sprintf("%s<a href='%s'>@%s</a>", jsKfTemp, kfUrlConfig.Value, kfNameConfig.Value), &tg.SendOptions{
				ReplyTo:   msgObj,
				ParseMode: tg.ModeHTML,
			})
		}

		//更新个人信息
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
			ctx.Send(gameLockedTemp, &tg.SendOptions{
				ReplyTo: msgObj,
			})
			return err
		}
		//大100、小100、单100、双100、大单100、大双100、小单100、小双100、6杀100(表示6点下注100)、对子100、顺子100、豹子100
		xz1Rx := regexp.MustCompile(`^(大|小|单|双|大单|大双|小单|小双|对子|顺子|豹子)(\d+)`)
		match = xz1Rx.FindStringSubmatch(txt)
		if match != nil {
			//过滤封盘
			qs, err := fnQsLocked(ctx)
			if err != nil {
				return err
			}

			xzType := match[1]
			xzMoney := match[2]
			xzMoneyF := mathutil.MustFloat(xzMoney)
			if user.Money < xzMoneyF {
				return ctx.Send(moneyBZTemp, &tg.SendOptions{
					ReplyTo: msgObj,
				})
			}
			configList, err := model.ConfigList()
			if err != nil {
				ctx.Send(jsKfTemp, &tg.SendOptions{
					ParseMode: tg.ModeHTML,
				})
				return err
			}
			stake, rate := xzInfoBytxTypeStr(xzType, configList)
			orderData := &model.Order{
				UserId:   user.ID,
				Username: user.Username,
				Nickname: user.Nickname,
				TgId:     user.TgId,
				Status:   0,
				Stake:    stake,
				Rate:     rate,
				QsId:     qs.ID,
				QsSn:     qs.Sn,
				Money:    xzMoneyF,
			}

			//保存数据
			tx := dao.Mdb.Begin()
			user.FreezMoney += xzMoneyF
			user.Money -= xzMoneyF
			user, err := model.EditUser(user)
			if err != nil {
				tx.Rollback()
				return err
			}
			order, err := model.AddOrder(orderData)
			if err != nil {
				tx.Rollback()
				return err
			}
			billData := &model.Bill{
				UserId:   user.ID,
				TgId:     user.TgId,
				Username: user.Username,
				Nickname: user.Nickname,
				Type:     3,
				ResId:    order.ID,
				Money:    xzMoneyF,
				Remark:   fmt.Sprintf(billRemarkTemp, user.Nickname, qs.Sn, xzType, xzMoneyF, order.Rate),
			}
			_, err = model.AddBill(billData)
			if err != nil {
				tx.Rollback()
				return err
			}
			tx.Commit()
			//发送下单通知
			beginDate := time.Unix(qs.BeginTime, 0).Format("15:04:05")
			endDate := time.Unix(qs.EndTime, 0).Format("15:04:05")
			orderList, err := model.GetOrderByQsIdAndStatus(qs.ID, 0)
			if err != nil {
				return nil
			}
			tzItemStr := ""
			for _, item := range *orderList {
				tzItemStr += fmt.Sprintf(xzGameStr2, xzInfoBytxTypeInt(item.Stake), item.Money, item.Rate)
			}
			replyMsg := fmt.Sprintf(xzGameStr, qs.Sn, beginDate, endDate, tzItemStr, user.Nickname, user.Money)
			return ctx.Send(replyMsg, &tg.SendOptions{
				ReplyTo:     msgObj,
				ParseMode:   tg.ModeHTML,
				ReplyMarkup: &tg.ReplyMarkup{InlineKeyboard: fnGroupInKeyBoard(qs.Sn)},
			})
		}
		xz2Rx := regexp.MustCompile(`^(\d+)杀(\d+)`)
		match = xz2Rx.FindStringSubmatch(txt)
		if match != nil {
			//过滤封盘
			qs, err := fnQsLocked(ctx)
			if err != nil {
				return err
			}

			xzType := match[1]
			xzTypeI := mathutil.MustInt(xzType)
			xzMoney := match[2]
			xzMoneyF := mathutil.MustFloat(xzMoney)
			if user.Money < xzMoneyF {
				return ctx.Send(moneyBZTemp, &tg.SendOptions{
					ReplyTo: msgObj,
				})
			}
			if xzTypeI < 3 || xzTypeI > 18 {
				return ctx.Send("点杀在3-18之间", &tg.SendOptions{
					ReplyTo: msgObj,
				})
			}

			configList, err := model.ConfigList()
			if err != nil {
				return ctx.Send(jsKfTemp, &tg.SendOptions{
					ParseMode: tg.ModeHTML,
				})
			}

			stake, rate := xzInfoBytxTypeStr(xzType, configList)
			orderData := &model.Order{
				UserId:   user.ID,
				TgId:     user.TgId,
				Username: user.Username,
				Nickname: user.Nickname,
				Rate:     rate,
				Stake:    stake,
				Status:   0,
				QsId:     qs.ID,
				QsSn:     qs.Sn,
				Money:    xzMoneyF,
			}
			//保存数据
			tx := dao.Mdb.Begin()
			user.FreezMoney += xzMoneyF
			user.Money -= xzMoneyF
			user, err := model.EditUser(user)
			if err != nil {
				tx.Rollback()
				return err
			}
			order, err := model.AddOrder(orderData)
			if err != nil {
				tx.Rollback()
				return err
			}
			billData := &model.Bill{
				UserId:   user.ID,
				TgId:     user.TgId,
				Username: user.Username,
				Nickname: user.Nickname,
				Type:     3,
				ResId:    order.ID,
				Money:    xzMoneyF,
				Remark:   fmt.Sprintf(billRemarkTemp, user.Nickname, qs.Sn, xzType, xzMoneyF, order.Rate),
			}
			_, err = model.AddBill(billData)
			if err != nil {
				tx.Rollback()
				return err
			}
			tx.Commit()

			//发送下单通知
			beginDate := time.Unix(qs.BeginTime, 0).Format("15:04:05")
			endDate := time.Unix(qs.EndTime, 0).Format("15:04:05")
			orderList, err := model.GetOrderByQsIdAndStatus(qs.ID, 0)
			if err != nil {
				return nil
			}
			tzItemStr := ""
			for _, item := range *orderList {
				tzItemStr += fmt.Sprintf(xzGameStr2, xzInfoBytxTypeInt(item.Stake), item.Money, item.Rate)
			}
			replyMsg := fmt.Sprintf(xzGameStr, qs.Sn, beginDate, endDate, tzItemStr, user.Nickname, user.Money)
			return ctx.Send(replyMsg, &tg.SendOptions{
				ReplyTo:     msgObj,
				ParseMode:   tg.ModeHTML,
				ReplyMarkup: &tg.ReplyMarkup{InlineKeyboard: fnGroupInKeyBoard(qs.Sn)},
			})
		}
		//删除无关消息
		//todo
		ctx.Delete()
	}

	if ctx.Message().Private() {
		//个人聊天
		return ctx.Send(jsKfTemp, &tg.SendOptions{
			ParseMode: tg.ModeHTML,
		})
	}
	return nil
}

func fnGroupInKeyBoard(qsSn string) [][]tg.InlineButton {
	kfUrlConfig, _ := model.GetConfigByName("kf_url")
	userMoneyBtn := tg.InlineButton{
		Text:   "当前余额",
		Unique: "userMoney",
	}
	Bots.Handle(&userMoneyBtn, fnUserMoney)
	curStakeBtn := tg.InlineButton{
		Text:   "当前投注",
		Unique: "curStake",
		Data:   qsSn,
	}
	Bots.Handle(&curStakeBtn, fnCurStake)

	historyStakeBtn := tg.InlineButton{
		Text:   "最近投注",
		Unique: "historyStake",
	}
	Bots.Handle(&historyStakeBtn, fnHistoryStake)

	billBtn := tg.InlineButton{
		Text:   "账单记录",
		Unique: "bill",
	}
	Bots.Handle(&billBtn, fnBillBtn)
	rechargeBtn := tg.InlineButton{
		Text:   fmt.Sprintf("我要充值"),
		Unique: "recharge",
		URL:    kfUrlConfig.Value,
	}
	withdrawalBtn := tg.InlineButton{
		Text:   fmt.Sprintf("我要提现"),
		Unique: "recharge",
		URL:    kfUrlConfig.Value,
	}
	btns := [][]tg.InlineButton{{userMoneyBtn, curStakeBtn}, {historyStakeBtn, billBtn}, {rechargeBtn, withdrawalBtn}}
	return btns
}

// ========================================================
// ======================群信息==========================
// ========================================================
// 账单记录
func fnBillBtn(ctx tg.Context) error {
	msg := "最近无账单记录"
	//更新个人信息
	sender := ctx.Callback().Sender
	list, err := model.GetBillByTgId(strutil.MustString(sender.ID), 5)
	if err != nil {
		return ctx.Respond(&tg.CallbackResponse{
			CallbackID: ctx.Callback().ID,
			Text:       msg,
			ShowAlert:  true,
		})
	}
	if len(*list) > 0 {
		msg = ""
		for _, item := range *list {
			msg += time.Unix(item.CreateTime, 0).Format("2006-01-02")
			if item.Type == 1 {
				msg += fmt.Sprintf("\t充值：+%.2f", item.Money)
			}
			if item.Type == 2 {
				msg += fmt.Sprintf("\t提现：-%.2f", item.Money)
			}
			if item.Type == 3 {
				msg += fmt.Sprintf("\t下注：-%.2f", item.Money)
			}
			if item.Type == 4 {
				msg += fmt.Sprintf("\t中奖：+%.2f", item.Money)
			}
			msg += "\n"
		}
	}
	return ctx.Respond(&tg.CallbackResponse{
		CallbackID: ctx.Callback().ID,
		Text:       msg,
		ShowAlert:  true,
	})

}

// 最近投注
func fnHistoryStake(ctx tg.Context) error {
	msg := "最近无投注记录"
	//更新个人信息
	sender := ctx.Callback().Sender
	list, err := model.GetOrderByTgId(strutil.MustString(sender.ID), 5)
	if err != nil {
		return ctx.Respond(&tg.CallbackResponse{
			CallbackID: ctx.Callback().ID,
			Text:       msg,
			ShowAlert:  true,
		})
	}
	if len(*list) > 0 {
		msg = ""
		for _, item := range *list {
			stake := xzInfoBytxTypeInt(item.Stake)
			money := "中奖：-"
			if item.Status == 1 {
				money = fmt.Sprintf("中奖：+%.2f", item.ResMoney)
			}

			if item.Status == 2 {
				money = fmt.Sprintf("中奖：-%.2f", item.ResMoney)
			}
			msg += fmt.Sprintf("%s-%.2f(赔率 1:%.2f)\t%s\n", stake, item.Money, item.Rate, money)
		}
	}
	return ctx.Respond(&tg.CallbackResponse{
		CallbackID: ctx.Callback().ID,
		Text:       msg,
		ShowAlert:  true,
	})
}

// 当前投注
func fnCurStake(ctx tg.Context) error {

	msg := "当前无投注记录"
	//更新个人信息
	sender := ctx.Callback().Sender
	qsSn := strutil.MustString(ctx.Data())
	if qsSn == "" {
		return ctx.Respond(&tg.CallbackResponse{
			CallbackID: ctx.Callback().ID,
			Text:       msg,
			ShowAlert:  true,
		})
	}

	list, err := model.GetOrderByQsSnAndTgId(qsSn, strutil.MustString(sender.ID))
	if err != nil {
		return ctx.Respond(&tg.CallbackResponse{
			CallbackID: ctx.Callback().ID,
			Text:       msg,
			ShowAlert:  true,
		})
	}
	if len(*list) > 0 {
		msg = ""
		for _, item := range *list {
			stake := xzInfoBytxTypeInt(item.Stake)
			msg += fmt.Sprintf("%s - %.2f （赔率 1:%.2f）\n", stake, item.Money, item.Rate)
		}
	}
	return ctx.Respond(&tg.CallbackResponse{
		CallbackID: ctx.Callback().ID,
		Text:       msg,
		ShowAlert:  true,
	})
}

// 当前余额
func fnUserMoney(ctx tg.Context) error {
	//更新个人信息
	sender := ctx.Callback().Sender
	user, err := model.GetUserInfoByTgId(strutil.MustString(sender.ID))
	if err != nil {
		return err
	}
	return ctx.Respond(&tg.CallbackResponse{
		CallbackID: ctx.Callback().ID,
		Text:       fmt.Sprintf("【账户ID】：%s\n【账户昵称】：%s\n【账户余额】：%.2f", user.TgId, user.Nickname, user.Money),
		ShowAlert:  true,
	})
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
	msgTemp += "【账户ID】：<code>%s</code>\n【账户昵称】：<b>%s</b>\n【账户余额】：<span class='tg-spoiler'>%.2f</span>\n"
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
