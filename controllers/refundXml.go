package controllers

import "encoding/xml"

//申请退款请求报文结构体
type RefundRequestXml struct {
	XMLName  					xml.Name 				`xml:"xml"`
	Appid 						string 					`xml:"appid"`					//公众账号ID
	Mch_id 						string 					`xml:"mch_id"`					//商户号
	Nonce_str					string					`xml:"nonce_str"`				//随机字符串
	Out_refund_no				string					`xml:"out_refund_no"`			//退款订单
	Refund_fee					string					`xml:"refund_fee"`			//本次退款金额
	Refund_fee_type				string					`xml:"refund_fee_type"`		//退款货币种类
	Out_trade_no				string					`xml:"out_trade_no"`			//商户订单号(原交易订单)
	Total_fee					string					`xml:"total_fee"`				//原交易订单金额
	Refund_desc					string					`xml:"refund_desc"`			//退款原因描述
	Notify_url					string					`xml:"notify_url"`			//退款结果通知url
	Sign						string					`xml:"sign"`					//签名
}

//申请退款返回报文结构体
type RefundResponseXml struct {
	XMLName  					xml.Name 				`xml:"xml"`

}

//退款结果异步通知
type RefundNotifyResult struct {
	XMLName						xml.Name				`xml:"xml"`

}

//退款异步通知-回复
type RefundNotifyResponse struct {
	XMLName 					xml.Name
	Return_code					string					`xml:"return_code"`			//返回状态码
	Return_msg					string					`xml:"return_msg"`			//返回信息
}







