package main

import (
	_ "unifiedpay/startup/globalSessions"
	_ "unifiedpay/startup/memory"
	_ "unifiedpay/startup/logger"
	_ "unifiedpay/routers"
	_ "unifiedpay/startup/db"
	"github.com/astaxie/beego"
)

func main() {
	beego.Run()
}

