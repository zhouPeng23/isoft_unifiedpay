package controllers

import (
	"strconv"
	"unifiedpay/models"
	"time"
	"math/rand"
	"fmt"
	"github.com/astaxie/beego/logs"
)

//下单发送https请求-对接微信支付
func (this *MainController)WeChatPay() error{
	//界面接收的参数
	productId := this.GetString("ProductId")
	productDesc := this.GetString("ProductDesc")
	transAmount, _ := strconv.Atoi(this.GetString("TransAmount"))
	transCurrCode := this.GetString("TransCurrCode")
	//生成订单随机数
	rand.Seed(time.Now().UnixNano())
	random := rand.Intn(999999999)
	strRandom := "000000000" + fmt.Sprintf("%d", random)
	strRandom = strRandom[len(strRandom)-9:len(strRandom)]
	now := time.Now().Format("20060102150405")
	//组装订单
	order := models.Order{}
	order.OrderId = now + strRandom
	order.PayType = "微信支付"
	order.TransType = "SALE"
	order.MerchantNo = "43119023467211"
	order.ProductId = productId
	order.ProductDesc = productDesc
	order.TransTime = now
	order.TransAmount = int64(transAmount)
	order.TransCurrCode = transCurrCode
	//入库
	e := order.Pay(order)
	if e != nil {
		logs.Error(e)
		return e
	}
	//发送微信下单请求


	return nil
}


