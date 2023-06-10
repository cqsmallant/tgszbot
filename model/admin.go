package model

import (
	"ant/utils/dao"
)

type AuthRule struct {
	title uint64 `gorm:"column:title" json:"title"`
	BaseModel
}

func (bill *AuthRule) TableName() string {
	return "fa_auth_rule"
}

func GetAuthRuleById(id uint64) (*AuthRule, error) {
	info := new(AuthRule)
	err := dao.Mdb.Model(info).Where("id = ?", id).Find(info).Error
	return info, err
}
