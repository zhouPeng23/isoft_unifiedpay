package models

import (
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/logs"
	"errors"
	"strings"
)

type Order struct {
	Id int64
	OrderId string `orm:"unique"`
	OrgOrderId string
	TransType string
	MerchantNo string
	ProductId string
	ProductDesc string
	TransTime string
	TransAmount int64
	TransCurrCode string
	TransStatus string
	RefundReason string
	RefundedAmount int64
}


//下单
func (this *Order)Pay(order Order) error {
	//1.支付参数验证
	validation := this.PayParamValidation(order)
	if validation!=nil {
		return validation
	}
	//2.入库
	o := orm.NewOrm()
	_, e := o.Insert(&order)
	if e != nil {
		logs.Error("支付订单入库失败:",order.OrderId,"失败原因:",e)
		return errors.New("支付订单入库失败！")
	}
	logs.Info("支付订单入库成功",order.OrderId)
	return nil
}


//退货
func (this *Order)Refund(order Order) error {
	//1.退货参数验证
	validation := this.RefundParamValidation(order)
	if validation!=nil {
		return validation
	}
	//2.查询原交易订单
	o := orm.NewOrm()
	originalOrder := Order{}
	originalOrder.OrderId = order.OrgOrderId
	err := o.Read(&originalOrder, "OrderId")
	if err != nil {
		return errors.New("原交易订单查询失败！")
	}
	//3.判断退货金额是否大于可退金额
	returnableAmount := originalOrder.TransAmount - originalOrder.RefundedAmount
	refundAmount := order.TransAmount
	if refundAmount > returnableAmount{
		return errors.New("退货金额大于可退金额！")
	}
	//4.入库
	_, e := o.Insert(&order)
	if e != nil {
		logs.Error("退货订单入库失败:",order.OrderId,"原交易订单:",order.OrgOrderId,"失败原因:",e)
		return errors.New("退货订单入库失败")
	}
	//5.设置原交易已退金额
	originalOrder.RefundedAmount = originalOrder.RefundedAmount + refundAmount
	o.Update(&originalOrder,"RefundedAmount")
	return nil
}


//支付订单参数验证
func (this *Order)PayParamValidation(order Order) error {
	if len(strings.TrimSpace(order.OrderId)) == 0 {
		return errors.New("交易订单为空！")
	}
	if order.TransType != "SALE" {
		return errors.New("交易类型不正确！")
	}
	if len(strings.TrimSpace(order.MerchantNo)) == 0 {
		return errors.New("商户号为空！")
	}
	if len(strings.TrimSpace(order.ProductId)) == 0 {
		return errors.New("下单的产品代码为空！")
	}
	if len(strings.TrimSpace(order.ProductDesc)) == 0 {
		return errors.New("下单的产品描述为空！")
	}
	if len(strings.TrimSpace(order.TransTime)) == 0 {
		return errors.New("交易时间为空！")
	}
	if order.TransAmount<=0 {
		return errors.New("交易金额不正确！")
	}
	if order.TransCurrCode != "156" {
		return errors.New("交易币种不正确，目前只支持RMB！")
	}
	return nil
}


//退货参数验证
func (this *Order)RefundParamValidation(order Order) error {
	if len(strings.TrimSpace(order.OrderId)) == 0 {
		return errors.New("退货订单为空！")
	}
	if len(strings.TrimSpace(order.OrgOrderId)) == 0 {
		return errors.New("原交易订单为空！")
	}
	if order.TransType!="REFUND" {
		return errors.New("交易类型不正确！")
	}
	if len(strings.TrimSpace(order.MerchantNo)) == 0 {
		return errors.New("商户号为空！")
	}
	if len(strings.TrimSpace(order.ProductId)) == 0 {
		return errors.New("退货的产品代码为空！")
	}
	if len(strings.TrimSpace(order.ProductDesc)) == 0 {
		return errors.New("退货的产品描述为空！")
	}
	if len(strings.TrimSpace(order.TransTime)) == 0 {
		return errors.New("交易时间为空！")
	}
	if order.TransAmount<=0 {
		return errors.New("交易金额不正确！")
	}
	if order.TransCurrCode != "156" {
		return errors.New("交易币种不正确，目前只支持RMB！")
	}
	if len(strings.TrimSpace(order.RefundReason)) == 0 {
		return errors.New("退货原因为空！")
	}
	return nil
}