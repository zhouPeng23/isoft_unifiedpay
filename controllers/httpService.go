package controllers

import (
	"unifiedpay/models"
	"time"
	"math/rand"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/httplib"
	"crypto/tls"
)

//下单发送https请求-对接微信支付
func (this *MainController)Pay(){
	go this.WeChatPay()
}

func (this *MainController)WeChatPay() error {
	//界面接收的参数
	//productId := this.GetString("ProductId")
	//productDesc := this.GetString("ProductDesc")
	//transAmount, _ := strconv.Atoi(this.GetString("TransAmount"))
	//transCurrCode := this.GetString("TransCurrCode")
	productId := "001256"
	productDesc := "苹果手机"
	transAmount := 88
	transCurrCode := "156"
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
	order.MerchantNo = beego.AppConfig.String("MerchantNo")
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
	req := httplib.Post(beego.AppConfig.String("WeChatPayUrl"))
	req.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	req.SetTimeout(60*time.Second,60*time.Second)
	req.XMLBody("")
	//获取返回消息
	var res interface{}
	json := req.ToJSON(&res)
	fmt.Printf("微信支付返货消息:%v\n",json)
	return nil
}


