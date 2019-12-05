package controllers

import (
	"unifiedpay/models"
	"time"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego"
	"sync"
	"math/rand"
	"github.com/astaxie/beego/httplib"
	"crypto/tls"
)

var (
	intRandom = 100000000
	lock sync.Mutex
)


//支付
func (this *MainController)Pay(){
	go this.WeChatPay()
}


//下单发送https请求-对接微信支付
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
	now := time.Now().Format("20060102150405")
	//组装订单
	order := models.Order{}
	order.OrderId = now + queryUniqueRandom()
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


//生成订单随机数
func queryUniqueRandom() string {
	if intRandom == 999999999{
		intRandom = 100000000
	}
	lock.Lock()
	intRandom++
	lock.Unlock()
	//获取一个长度为9的唯一数字字符串（给自己人看的）
	intStr9 := "000000000" + fmt.Sprintf("%d",intRandom)
	intStr9 = intStr9[len(intStr9)-9:len(intStr9)]
	//获取一个长度为7的随机数（用于干扰别有用心者）
	rand.Seed(time.Now().UnixNano())
	random := rand.Intn(7777777)
	strRandom7 := "0000000" + fmt.Sprintf("%d", random)
	strRandom7 = strRandom7[len(strRandom7)-7:len(strRandom7)]
	return intStr9 + strRandom7
}
