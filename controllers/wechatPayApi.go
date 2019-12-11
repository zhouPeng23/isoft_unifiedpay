package controllers

import (
	"isoft_unifiedpay/models"
	"github.com/astaxie/beego/orm"
	"encoding/json"
	"github.com/astaxie/beego/logs"
	"fmt"
)

//订单查询
func (this *MainController)QueryOrder(){
	OrderId := this.GetString("OrderId")
	TransType := this.GetString("TransType")
	ProductDesc := this.GetString("ProductDesc")
	TransTime := this.GetString("TransTime")
	TransAmount := this.GetString("TransAmount")
	logs.Info(fmt.Sprintf("接口入参: OrderId=%v, TransType=%v, ProductDesc=%v, TransTime=%v, TransAmount=%v",OrderId,TransType,ProductDesc,TransTime,TransAmount))
	o := orm.NewOrm()
	order := models.Order{}
	order.OrderId = OrderId
	//order.TransType = TransType
	//order.ProductDesc = ProductDesc
	//order.TransTime = TransTime
	//amount, _ := strconv.Atoi(TransAmount)
	//order.TransAmount = int64(amount)
	order.OrderId = OrderId
	o.Read(&order,"OrderId")
	dataBytes, _ := json.Marshal(order)
	this.Data["json"] = string(dataBytes)
	this.ServeJSON()
}
