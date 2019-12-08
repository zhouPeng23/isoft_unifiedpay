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
	"strings"
)


//支付下单-控制器
func (this *MainController)Pay(){
	go this.WeChatPay()
}


//支付结果异步通知-控制器
func (this *MainController)PayNotifyResult() {
	go this.WeChatPayNofify()
}


//支付下单-具体处理方法
func (this *MainController)WeChatPay() (string,error) {
	code_url := ""//支付二维码
	o := orm.NewOrm()
	//o.Begin()
	//界面接收的参数
	//productId := this.GetString("ProductId")
	//productDesc := this.GetString("ProductDesc")
	//transAmount := this.GetString("TransAmount")
	//transCurrCode := this.GetString("TransCurrCode")
	productId := "001256"
	productDesc := "苹果手机"
	transAmount := "5888"
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
			//下单成功
			orderSuccess := models.Order{}
			orderSuccess.OrderId = order.OrderId
			o.Read(&orderSuccess, "OrderId")
			orderSuccess.OrderResultCode = resXml.Result_code
			orderSuccess.OrderResultDesc = "下单成功"
			orderSuccess.CodeUrl = resXml.Code_url
			o.Update(&orderSuccess)
			//这里设置真正的支付二维码
			code_url = resXml.Code_url
		}else {
			//下单失败
			orderFail := models.Order{}
			orderFail.OrderId = order.OrderId
			o.Read(&orderFail, "OrderId")
			orderFail.OrderResultCode = resXml.Result_code
			orderFail.OrderResultDesc = "下单失败"
			orderFail.OrderErrCode = resXml.Err_code
			orderFail.OrderErrCodeDes = resXml.Err_code_des
			o.Update(&orderFail)
			return code_url,errors.New(fmt.Sprintf("下单失败,失败原因:%v",resXml.Err_code_des))
		}
		o.Commit()
	}else {
		//通信都失败了，直接回滚
		//logs.Info("通信都失败了，直接回滚")
		//o.Rollback()
		return code_url,errors.New(fmt.Sprintf("通信标识:FAIL,失败原因:%v",resXml.Return_msg))
	}

	return code_url,nil
}


//支付结果异步通知-具体处理方法
func (this *MainController)WeChatPayNofify() {
	//获取支付结果异步通知
	logs.Info("支付结果异步通知上来了...")
	reqBody := this.Ctx.Input.RequestBody
	reqXml := NativePayNotifyResult{}
	e := xml.Unmarshal([]byte(reqBody), &reqXml)
	if e != nil {
		logs.Info("支付结果异步通知失败:报文转结构体失败")
	}else {
		logs.Info("支付结果异步通知转结构体成功")
	}

	//开始解析结构体
	logs.Info("开始解析结构体")
	if reqXml.Return_code == "SUCCESS" {
		o := orm.NewOrm()
		order := models.Order{}
		orderId := reqXml.Out_trade_no
		order.OrderId = orderId
		e := o.Read(&order, "OrderId")
		if e!=nil {
			logs.Info(fmt.Sprintf("原交易订单%v查询失败！",orderId))
		}
		order.PayResultCode = reqXml.Result_code
		if reqXml.Result_code=="SUCCESS" {
			//支付成功
			logs.Info("支付成功")
			order.PayResultDesc = "支付成功"
			order.TransactionId = reqXml.Transaction_id
			if len(strings.TrimSpace(reqXml.Cash_fee))>0 {//暂时认为是微信零钱支付
				order.WechatCash = "微信零钱支付"
			}
			if len(strings.TrimSpace(reqXml.Bank_type))>0 {//银行卡支付
				order.BankType = reqXml.Bank_type
				order.BankName = beego.AppConfig.String(reqXml.Bank_type)//bank.conf
			}
			order.PayTimeEnd = reqXml.Time_end
		}else {
			//支付失败
			logs.Info("支付失败")
			order.PayResultDesc = "支付失败"
			order.PayErrCode = reqXml.Err_code
			order.PayErrCodeDesc = reqXml.Err_code_des
		}
		o.Update(&order)

		//收到微信支付结果通知后，给一个成功应答
		resXml := NativePayNotifyResponse{}
		resXml.Return_code = "SUCCESS"
		resXml.Return_msg = "OK"
		resXmlStr, e := xml.Marshal(resXml)
		this.Data["xml"] = resXmlStr
		this.ServeXML()
	}else {
		logs.Info(fmt.Sprintf("返回状态码失败，失败原因:%v",reqXml.Return_msg))
	}
}






