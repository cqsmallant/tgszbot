package model

import (
	"ant/utils/dao"
)

type Config struct {
	ID    uint64 `gorm:"column:id"`
	Name  string `gorm:"column:name"`
	Value string `gorm:"column:value"`
}

func (user *Config) TableName() string {
	return "fa_config"
}

// GetUserInfoByTgId 通过客户信息
func ConfigList() ([]Config, error) {
	list := []Config{}
	err := dao.Mdb.Model(&Config{}).Find(&list).Error
	return list, err
}

// GetConfigByName 通过客户信息
func GetConfigByName(name string) (*Config, error) {
	info := new(Config)
	err := dao.Mdb.Model(info).Where(" name = ?", name).Find(info).Error
	return info, err
}
