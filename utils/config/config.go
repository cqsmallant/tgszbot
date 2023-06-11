package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

var (
	AppName           string
	AppUrl            string
	AppDeBug          bool
	HttpListen        string
	StaticPath        string
	RuntimePath       string
	MysqlDns          string
	MysqlTablePrefix  string
	MysqlMaxIdleConns int
	MysqlMaxOpenConns int
	MysqlMaxLifeTime  int

	RedisDns         string
	RedisDb          int
	RedisPwd         string
	RedisPooSize     int
	RedisMaxRetries  int
	RedisIdleTimeout int

	LogSavePath   string
	LogMaxSize    int
	LogMaxBackups int
	LogMaxAge     int
	TgBotToken    string
	TgGroupId     int64
	ApiProxy      string
	TrcToken      string
	QsStep        int

	QueueConcurrency   int
	QueueLevelCritical int
	QueueLevelDefault  int
	QueueLevelLow      int
)

func init() {
	viper.AddConfigPath("./")
	viper.SetConfigFile("app.env")
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	gwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	viper.GetInt("redis_db")

	AppName = viper.GetString("app_name")
	AppUrl = viper.GetString("app_url")
	AppDeBug = viper.GetBool("app_debug")
	HttpListen = viper.GetString("http_listen")
	StaticPath = viper.GetString("static_path")
	RuntimePath = fmt.Sprintf("%s%s", gwd, viper.GetString("runtime_path"))

	LogSavePath = fmt.Sprintf("%s%s", RuntimePath, viper.GetString("log_save_path"))
	LogMaxSize = viper.GetInt("log_max_size")

	TgBotToken = viper.GetString("tg_bot_token")
	TgGroupId = viper.GetInt64("tg_group_id")
	ApiProxy = viper.GetString("api_proxy")
	//充值地址
	TrcToken = viper.GetString("trc_token")
	QsStep = viper.GetInt("qs_step")

	//Mysql
	MysqlDns = fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		viper.GetString("mysql_user"),
		viper.GetString("mysql_passwd"),
		fmt.Sprintf(
			"%s:%s",
			viper.GetString("mysql_host"),
			viper.GetString("mysql_port")),
		viper.GetString("mysql_database"))
	MysqlTablePrefix = viper.GetString("mysql_table_prefix")
	MysqlMaxIdleConns = viper.GetInt("mysql_max_idle_conns")
	MysqlMaxOpenConns = viper.GetInt("mysql_max_open_conns")
	MysqlMaxLifeTime = viper.GetInt("mysql_max_life_time")

	//redis
	RedisDns = fmt.Sprintf(
		"%s:%s",
		viper.GetString("redis_host"),
		viper.GetString("redis_port"))
	RedisDb = viper.GetInt("redis_db")
	RedisPwd = viper.GetString("redis_passwd")
	RedisPooSize = viper.GetInt("redis_poo_size")
	RedisMaxRetries = viper.GetInt("redis_max_retries")
	RedisIdleTimeout = viper.GetInt("redis_idle_timeout")

	QueueConcurrency = viper.GetInt("queue_concurrency")
	QueueLevelCritical = viper.GetInt("queue_level_critical")
	QueueLevelDefault = viper.GetInt("queue_level_default")
	QueueLevelLow = viper.GetInt("queue_level_low")
}

func GetAppVersion() string {
	return "1.0.0"
}
