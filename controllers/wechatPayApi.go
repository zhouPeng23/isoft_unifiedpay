package controllers

import (
	"isoft_unifiedpay/models"
	"github.com/astaxie/beego/orm"
	"encoding/json"
	"github.com/astaxie/beego/logs"
	"fmt"
	"strconv"
)

//订单查询
func (this *MainController)QueryOrder(){
	logs.Info("调用订单查询接口...")
	OrderId := this.GetString("OrderId")
	TransType := this.GetString("TransType")
	ProductDesc := this.GetString("ProductDesc")
	TransTime := this.GetString("TransTime")
	TransAmount := this.GetString("TransAmount")
	logs.Info(fmt.Sprintf("接口入参: OrderId=%v, TransType=%v, ProductDesc=%v, TransTime=%v, TransAmount=%v",OrderId,TransType,ProductDesc,TransTime,TransAmount))
	o := orm.NewOrm()
	var orders []*models.Order;
	qs := o.QueryTable("Order").Limit(100)
	if len(OrderId)>0 {
		qs = qs.Filter("OrderId__istartswith",OrderId)
	}
	if len(TransType)>0 {
		qs = qs.Filter("TransType", TransType)
	}
	if len(ProductDesc)>0 {
		qs = qs.Filter("ProductDesc__icontains", ProductDesc)
	}
	if len(TransTime)>0 {
		qs = qs.Filter("TransTime__istartswith",TransTime)
	}
	if len(TransAmount)>0 {
		amount, _ := strconv.ParseFloat(TransAmount,64)//string 转 float
		qs = qs.Filter("TransAmount",int64(amount*100))
	}
	qs.OrderBy("-TransTime").All(&orders)
	dataBytes, _ := json.Marshal(orders)
	this.Data["json"] = string(dataBytes)
	this.ServeJSON()
}
