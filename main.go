package main

import (
	"ant/bootstrap"
	"ant/utils/config"

	"github.com/gookit/color"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			color.Error.Println("[Start Server Err!!!] ", err)
		}
	}()
	color.Green.Printf("%s\n", " .----------------.  .-----------------. .----------------. \n| .--------------. || .--------------. || .--------------. |\n| |      __      | || | ____  _____  | || |  _________   | |\n| |     /  \\     | || ||_   \\|_   _| | || | |  _   _  |  | |\n| |    / /\\ \\    | || |  |   \\ | |   | || | |_/ | | \\_|  | |\n| |   / ____ \\   | || |  | |\\ \\| |   | || |     | |      | |\n| | _/ /    \\ \\_ | || | _| |_\\   |_  | || |    _| |_     | |\n| ||____|  |____|| || ||_____|\\____| | || |   |_____|    | |\n| |              | || |              | || |              | |\n| '--------------' || '--------------' || '--------------' |\n '----------------'  '----------------'  '----------------' ")
	color.Infof("%s version(%s) Powered by %s %s\n", config.AppName, config.GetAppVersion(), config.AppUrl, "@sunnant")
	bootstrap.Start()
}
