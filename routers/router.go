package routers

import (
	"isoft_unifiedpay/controllers"
	"github.com/astaxie/beego"
)

//统一使用默认MainController
var mc = &controllers.MainController{}

//简便注册，为了在具体路由里展示方法，可以鼠标点击进入
func registRouter(rootpath string, c beego.ControllerInterface, callFunc func(), mappingMethods ...string) *beego.App {
	return beego.Router(rootpath, c, mappingMethods...)
}

func init() {
    beego.Router("/", mc)
	registRouter("/wechatPayApi/Pay", mc,mc.Pay,"post:Pay")
	registRouter("/wechatPayApi/Refund", mc,mc.Refund,"post:Refund")
	registRouter("/wechatPayApi/QueryOrder", mc,mc.QueryOrder,"post:QueryOrder")
	registRouter("/wechatPayApi/ShowLastedOrders", mc,mc.ShowLastedOrders,"post:ShowLastedOrders")

}
