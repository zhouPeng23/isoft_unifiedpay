package main

import (
	_ "isoft_unifiedpay/startup/globalSessions"
	_ "isoft_unifiedpay/startup/memory"
	_ "isoft_unifiedpay/startup/logger"
	_ "isoft_unifiedpay/routers"
	_ "isoft_unifiedpay/startup/db"
	"github.com/astaxie/beego"
)

func main() {
	beego.Run()
}

