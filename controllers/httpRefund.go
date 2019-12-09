package controllers

import (
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/logs"
	"fmt"
	"time"
	"unifiedpay/models"
	"github.com/astaxie/beego"
	"strconv"
	"github.com/astaxie/beego/httplib"
	"crypto/tls"
	"encoding/xml"
	"errors"
)

//退货请求-控制器
func (this *MainController) Refund() {
	go this.WeChatRefund()
}


//退款结果异步通知-控制器
func (this *MainController)RefundNotifyResult() {
	go this.WeChatRefundNofify()
}


//退货请求-具体处理方法
func (this *MainController) WeChatRefund() (string,error) {
	applyResult := "" //申请结果，给先给个默认值
	o := orm.NewOrm()
	//发起退款申请，接收参数
	//orgOrderId := this.GetString("OrgOrderId")
	//transAmount := this.GetString("TransAmount")
	//transCurrCode := this.GetString("TransCurrCode")
	//refundReason := this.GetString("RefundReason")
	orgOrderId := "201912080943371000000062859163"
	transAmount := "20"
	transCurrCode := "CNY"
	refundReason := "手机发热严重"
	logs.Info("退货请求上来了...")
	logs.Info(fmt.Sprintf("请求参数:orgOrderId=%v,transAmount=%v,transCurrCode=%v,refundReason=%v",orgOrderId,transAmount,transCurrCode,refundReason))
	now := time.Now().Format("20060102150405")

	//查询原交易，获取商品基本参数(主要获取商户ID和描述)
	orgOrder := models.Order{}
	orgOrder.OrderId = orgOrderId
	o.Read(&orgOrder,"OrderId")

	//组装退款订单
	order := models.Order{}
	order.OrderId = now + QueryUniqueRandom()
	order.OrgOrderId = orgOrderId
	order.TransType = "REFUND"
	order.MerchantNo = beego.AppConfig.String("WeChatPay_MerchantNo")
	order.ProductId = orgOrder.ProductId
	order.ProductDesc = orgOrder.ProductDesc
	order.TransTime = now
	amount, _ := strconv.Atoi(transAmount)
	order.TransAmount = int64(amount)
	order.TransCurrCode = transCurrCode
	order.RefundReason = refundReason

	//入库
	lock.Lock()
	e := order.Refund(o,order)
	lock.Unlock()
	if e != nil {
		logs.Error(e)
		return applyResult,e
	}else {
		logs.Info(fmt.Sprintf("退货订单%v入库成功",order.OrderId))
	}

	//发送微信申请退款请求
	logs.Info("发送申请退款请求...")
	req := httplib.Post(beego.AppConfig.String("WeChatPay_RefundApply"))
	req.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	req.SetTimeout(60*time.Second,60*time.Second)

	//组织xml报文
	reqXml := RefundRequestXml{}
	reqXml.Appid = beego.AppConfig.String("WeChatPay_Appid")
	reqXml.Mch_id = beego.AppConfig.String("WeChatPay_MerchantNo")
	reqXml.Nonce_str = "2ddd1a30ac87aa2db72f57a2375d8fec"
	reqXml.Out_refund_no = order.OrderId
	reqXml.Refund_fee = strconv.Itoa(int(order.TransAmount))
	reqXml.Refund_fee_type = order.TransCurrCode
	reqXml.Out_trade_no = order.OrgOrderId
	reqXml.Total_fee = strconv.Itoa(int(orgOrder.TransAmount))
	reqXml.Refund_desc = order.RefundReason
	reqXml.Notify_url = beego.AppConfig.String("WeChatPay_RefNotifyUrl")
	reqXml.Sign = "3CB01533B8C1EF103065174F50BCA002"

	//设置xml报文体
	reqXmlStr, e := xml.Marshal(reqXml)
	logs.Info("设置xml报文体:%v",string(reqXmlStr))
	req.XMLBody(string(reqXmlStr))

	//获取返回消息、转为结构体
	logs.Info("接收返回报文...")
	resXmlStr, e := req.String()
	logs.Info(fmt.Sprintf("收到报文:%v",resXmlStr))
	resXml := RefundResponseXml{}
	e = xml.Unmarshal([]byte(resXmlStr), &resXml)
	if e != nil {
		return applyResult,errors.New(fmt.Sprintf("转换返回报文为结构体失败,失败原因:%v",e.Error()))
	}else {
		logs.Info("转换返回报文为结构体成功")
	}

	//开始解析结构体
	logs.Info("开始解析结构体...")
	if resXml.Return_code=="SUCCESS" {
		//通信成功，数据都入库
		if resXml.Result_code=="SUCCESS" {
			//退款申请成功
			orderSuccess := models.Order{}
			orderSuccess.OrderId = order.OrderId
			o.Read(&orderSuccess, "OrderId")
			orderSuccess.RefundReqResultCode = resXml.Result_code
			orderSuccess.RefundReqResultDesc = "退款申请成功"
			o.Update(&orderSuccess)
			applyResult = "退款申请成功"
		}else {
			//退款申请失败
			orderFail := models.Order{}
			orderFail.OrderId = order.OrderId
			o.Read(&orderFail, "OrderId")
			orderFail.RefundReqResultCode = resXml.Result_code
			orderFail.RefundReqResultDesc = "退款申请失败"
			orderFail.RefundReqErrCode = resXml.Err_code
			orderFail.RefundReqErrCodeDesc = resXml.Err_code_des
			o.Update(&orderFail)
			return applyResult,errors.New(fmt.Sprintf("退款申请失败,失败原因:%v",resXml.Err_code_des))
		}
	}else {
		return applyResult,errors.New(fmt.Sprintf("通信标识:FAIL,失败原因:%v",resXml.Return_msg))
	}

	return applyResult,nil
}


//退款结果异步通知-具体处理方法
func (this *MainController)WeChatRefundNofify() {
	//获取支付结果异步通知
	logs.Info("退款结果异步通知上来了...")
	reqBody := this.Ctx.Input.RequestBody
	reqXml := RefundNotifyResult{}
	e := xml.Unmarshal([]byte(reqBody), &reqXml)
	if e != nil {
		logs.Info("退款结果异步通知失败:报文转结构体失败")
	}else {
		logs.Info("退款结果异步通知转结构体成功")
	}

	//开始解析结构体
	logs.Info("开始解析结构体")
	if reqXml.Return_code == "SUCCESS" {
		//加密信息解密
		reqInfo := reqXml.Req_info
		reqInfo = DecodeRefundNotifyReqInfo(reqInfo)
		reqInfoXml := RefInfo{}
		xml.Unmarshal([]byte(reqInfo), &reqInfoXml)
		//更新退款结果信息入库
		o := orm.NewOrm()
		order := models.Order{}
		orderId := reqInfoXml.Out_refund_no
		order.OrderId = orderId
		e := o.Read(&order, "OrderId")
		if e!=nil {
			logs.Info(fmt.Sprintf("退款订单%v查询失败！",orderId))
		}
		if reqInfoXml.Refund_status=="SUCCESS" {
			//退款成功
			order.RefundStatus = "退款成功"
			order.RefundSuccessTime = reqInfoXml.Success_time
		}else {
			//退款失败
			order.RefundStatus = "退款失败"
			//退款失败，需要更新原交易已退金额字段
			orgOrder := models.Order{}
			orgOrder.OrderId = order.OrgOrderId
			o.Read(&orgOrder,"OrderId")
			orgOrder.RefundedAmount = orgOrder.RefundedAmount - order.TransAmount
			o.Update(&orgOrder)
		}
		o.Update(&order)

		//收到微信退款结果通知后，给一个成功应答
		resXml := RefundNotifyResponse{}
		resXml.Return_code = "SUCCESS"
		resXml.Return_msg = "OK"
		resXmlStr, e := xml.Marshal(resXml)
		this.Data["xml"] = resXmlStr
		this.ServeXML()
	}else {
		logs.Info(fmt.Sprintf("返回状态码失败，失败原因:%v",reqXml.Return_msg))
	}
}




