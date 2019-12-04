package logger

import (
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego"
)

func init()  {
	var logDir string
	runmode := beego.AppConfig.String("runmode")
	if runmode!="prod" {
		logDir = "../../unifiedpay_log"
	}else {
		logDir = "../../unifiedpay_log"//项目部署的同级目录下放置log
	}

	logs.EnableFuncCallDepth(true)
	logs.SetLogger(logs.AdapterConsole)
	logs.SetLogger(logs.AdapterMultiFile,`{"filename":"`+logDir+`unifiedpay.log","maxdays":90,"separate":["emergency", "alert", "critical", "error", "warning", "notice", "info", "debug"]}`)
	logs.Async(1e3)
}


