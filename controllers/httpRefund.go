package controllers

import (
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/logs"
	"fmt"
	"time"
	"unifiedpay/models"
	"github.com/astaxie/beego"
	"strconv"
)

//退货请求-控制器
func (this *MainController) Refund() {
	go this.WeChatRefund()
}


//退货请求-具体处理方法
func (this *MainController) WeChatRefund() (string,error) {
	var applyResult string//申请结果
	o := orm.NewOrm()
	o.Begin()
	//发起退款申请，接收参数
	//orgOrderId := this.GetString("OrgOrderId")
	//transAmount := this.GetString("TransAmount")
	//transCurrCode := this.GetString("TransCurrCode")
	//refundReason := this.GetString("RefundReason")
	orgOrderId := "201912071835301000000001234567"
	transAmount := "66"
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
	e := order.Refund(o,order)
	if e != nil {
		logs.Error(e)
		return applyResult,e
	}else {
		logs.Info(fmt.Sprintf("退货订单%v入库成功",order.OrderId))
	}

	//发送微信申请退款请求


	return applyResult,nil
}



