package dao

import (
	"ant/utils/config"
	"ant/utils/log"
	"time"

	"github.com/gookit/color"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var Mdb *gorm.DB

// MysqlInit 数据库初始化
func MysqlInit() {
	var err error
	Mdb, err = gorm.Open(mysql.Open(config.MysqlDns), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   config.MysqlTablePrefix,
			SingularTable: true,
		},
		Logger: logger.Default.LogMode(logger.Error),
	})
	if err != nil {
		panic(err)
	}
	if config.AppDeBug {
		Mdb = Mdb.Debug()
	}
	sqlDB, err := Mdb.DB()
	if err != nil {
		color.Red.Printf("[store_db] mysql get DB,err=%s\n", err)
		panic(err)
	}
	sqlDB.SetMaxIdleConns(config.MysqlMaxIdleConns)
	sqlDB.SetMaxOpenConns(config.MysqlMaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Hour * time.Duration(config.MysqlMaxLifeTime))
	err = sqlDB.Ping()
	if err != nil {
		color.Red.Printf("[store_db] mysql connDB err:%s", err.Error())
		panic(err)
	}
	log.Sugar.Debug("[store_db] mysql connDB success")
}
