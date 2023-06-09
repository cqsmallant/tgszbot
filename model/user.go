package model

import (
	"ant/utils/dao"
	"time"
)

type User struct {
	TgId       string  `gorm:"column:tg_id"`
	GroupId    int     `gorm:"column:group_id"`
	Username   string  `gorm:"column:username"`
	Nickname   string  `gorm:"column:nickname"`
	Money      float64 `gorm:"column:money"`
	FreezMoney float64 `gorm:"column:freez_money"`
	Status     string  `gorm:"column:status"`
	BaseModel
}

func (user *User) TableName() string {
	return "fa_user"
}

// GetUserInfoByTgId 通过客户信息
func GetUserInfoByTgId(tgId string) (*User, error) {
	user := &User{}
	err := dao.Mdb.Model(&User{}).Limit(1).Find(user, "tg_id = ?", tgId).Error
	return user, err
}

func AddUser(user *User) (*User, error) {
	user.CreateTime = time.Now().Unix()
	err := dao.Mdb.Model(user).Create(user).Error
	return user, err
}

func EditUser(user *User) (*User, error) {
	user.UpdateTime = time.Now().Unix()
	err := dao.Mdb.Model(user).Where("tg_id=?", user.TgId).Updates(user).Error
	return user, err
}
