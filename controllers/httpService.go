package controllers

import (
	"unifiedpay/models"
	"time"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/httplib"
	"crypto/tls"
	"encoding/xml"
	"github.com/astaxie/beego/orm"
	"strconv"
)


//支付
func (this *MainController)Pay(){
	go this.WeChatPay()
}


//下单发送https请求-对接微信支付
func (this *MainController)WeChatPay() error {
	o := orm.NewOrm()
	//界面接收的参数
	//productId := this.GetString("ProductId")
	//productDesc := this.GetString("ProductDesc")
	//transAmount, _ := strconv.Atoi(this.GetString("TransAmount"))
	//transCurrCode := this.GetString("TransCurrCode")
	productId := "001256"
	productDesc := "苹果手机"
	transAmount := "88"
	transCurrCode := "CNY"
	logs.Info("微信扫码支付请求上来了...")
	logs.Info(fmt.Sprintf("请求参数:productId=%v,productDesc=%v,transAmount=%v,transCurrCode=%v",productId,productDesc,transAmount,transCurrCode))
	now := time.Now().Format("20060102150405")

	//组装订单
	order := models.Order{}
	order.OrderId = now + QueryUniqueRandom()
	order.PayType = "微信支付"
	order.TransType = "SALE"
	order.MerchantNo = beego.AppConfig.String("WeChatPay_MerchantNo")
	order.ProductId = productId
	order.ProductDesc = productDesc
	order.TransTime = now
	amount, _ := strconv.Atoi(transAmount)
	order.TransAmount = int64(amount)
	order.TransCurrCode = transCurrCode

	//入库
	logs.Info("订单开始入库...")
	e := order.Pay(o,order)
	if e != nil {
		logs.Error(e)
		return e
	}else {
		logs.Info(fmt.Sprintf("订单%v入库成功!",order.OrderId))
	}

	//发送微信下单请求
	logs.Info("发送微信下单请求...")
	req := httplib.Post(beego.AppConfig.String("WeChatPay_ReqUrl"))
	req.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	req.SetTimeout(60*time.Second,60*time.Second)

	//组织xml报文
	reqXml := NativeRequestXml{}
	reqXml.Appid = beego.AppConfig.String("WeChatPay_Appid")
	reqXml.Mch_id = order.MerchantNo
	reqXml.Out_trade_no = order.OrderId
	reqXml.Trade_type = beego.AppConfig.String("WeChatPay_TradeType")
	reqXml.Fee_type = order.TransCurrCode
	reqXml.Total_fee = strconv.Itoa(int(order.TransAmount))
	reqXml.Body = order.ProductDesc
	reqXml.Spbill_create_ip = beego.AppConfig.String("WeChatPay_SpbillCreateIp")
	reqXml.Notify_url = beego.AppConfig.String("WeChatPay_NotifyUrl")
	reqXml.Nonce_str = "1add1a30ac87aa2db72f57a2375d8fec"
	reqXml.Sign = "0CB01533B8C1EF103065174F50BCA001"

	//设置xml报文体
	xmlStr, e := xml.Marshal(reqXml)
	logs.Info("设置xml报文体:%v",string(xmlStr))
	req.XMLBody(string(xmlStr))

	//获取返回消息
	logs.Info("获取返回消息...")
	response, e := req.String()
	logs.Info(fmt.Sprintf("收到应答:%v",response))

	return nil
}



