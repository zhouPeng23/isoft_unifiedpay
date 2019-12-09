package routers

import (
	"isoft_unifiedpay/controllers"
	"github.com/astaxie/beego"
)

func init() {
    beego.Router("/", &controllers.MainController{})

}
