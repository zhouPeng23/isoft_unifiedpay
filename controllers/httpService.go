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
	"errors"
)


//支付
func (this *MainController)Pay(){
	go this.WeChatPay()
}


//下单发送https请求-对接微信支付
func (this *MainController)WeChatPay() (string,error) {
	code_url := ""//支付二维码
	o := orm.NewOrm()
	o.Begin()
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
	order.PayStyle = "微信支付"
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
		return code_url,e
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
	reqXmlStr, e := xml.Marshal(reqXml)
	logs.Info("设置xml报文体:%v",string(reqXmlStr))
	req.XMLBody(string(reqXmlStr))

	//获取返回消息、转为结构体
	logs.Info("接收返回报文...")
	resXmlStr, e := req.String()
	logs.Info(fmt.Sprintf("收到报文:%v",resXmlStr))
	resXml := NativeResponseXml{}
	e = xml.Unmarshal([]byte(resXmlStr), &resXml)
	if e != nil {
		return code_url,errors.New(fmt.Sprintf("转换返回报文为结构体失败,失败原因:%v",e.Error()))
	}else {
		logs.Info("转换返回报文为结构体成功")
	}

	//开始解析结构体
	logs.Info("开始解析结构体...")
	if resXml.Return_code=="SUCCESS" {
		//通信成功，则不管用户后面是否支付成功，数据都入库
		if resXml.Result_code=="SUCCESS" {
			//获取付款二维码
			orderSuccess := models.Order{}
			orderSuccess.OrderId = order.OrderId
			o.Read(&orderSuccess, "OrderId")
			orderSuccess.OrderResultCode = resXml.Return_code
			orderSuccess.OrderResultDesc = resXml.Return_msg
			orderSuccess.CodeUrl = resXml.Code_url
			o.Update(&orderSuccess,"OrderResultCode","OrderResultDesc","CodeUrl")
			code_url = resXml.Code_url //这里设置真正的支付二维码
		}else {
			//查询订单，给订单设置"下单失败"原因
			orderFail := models.Order{}
			orderFail.OrderId = order.OrderId
			o.Read(&orderFail, "OrderId")
			orderFail.OrderResultCode = resXml.Return_code
			orderFail.OrderResultDesc = resXml.Return_msg
			o.Update(&orderFail,"OrderResultCode","OrderResultDesc")
		}
		o.Commit()
	}else {
		//通信都失败了，直接回滚
		logs.Info("通信都失败了，直接回滚")
		o.Rollback()
		return code_url,errors.New(fmt.Sprintf("通信标识:FAIL,失败原因:%v",resXml.Return_msg))
	}

	return code_url,nil
}



