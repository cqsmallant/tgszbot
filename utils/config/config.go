package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

var (
	AppName       string
	AppUrl        string
	AppDeBug      bool
	HttpListen    string
	RuntimePath   string
	MysqlDns      string
	LogSavePath   string
	LogMaxSize    int
	LogMaxBackups int
	LogMaxAge     int
	StaticPath    string
	TgBotToken    string
	TgManage      int64
	ApiProxy      string
	TrcToken      string
	QsStep        int
)

func init() {
	viper.AddConfigPath("./")
	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	gwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	AppName = viper.GetString("app_name")
	AppUrl = viper.GetString("app_url")
	AppDeBug = viper.GetBool("app_debug")
	HttpListen = viper.GetString("http_listen")
	StaticPath = viper.GetString("static_path")
	RuntimePath = fmt.Sprintf("%s%s", gwd, viper.GetString("runtime_path"))
	LogSavePath = fmt.Sprintf("%s%s", RuntimePath, viper.GetString("log_save_path"))
	LogMaxSize = viper.GetInt("log_max_size")
	TgBotToken = viper.GetString("tg_bot_token")
	TgManage = viper.GetInt64("tg_manage")
	ApiProxy = viper.GetString("api_proxy")
	TrcToken = viper.GetString("trc_token")
	QsStep = viper.GetInt("qs_step")
	MysqlDns = fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		viper.GetString("mysql_user"),
		viper.GetString("mysql_passwd"),
		fmt.Sprintf(
			"%s:%s",
			viper.GetString("mysql_host"),
			viper.GetString("mysql_port")),
		viper.GetString("mysql_database"))

}

func GetAppVersion() string {
	return "1.0.0"
}
