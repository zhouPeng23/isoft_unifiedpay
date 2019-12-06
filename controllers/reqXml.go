package controllers

import "encoding/xml"

//微信支付（Native支付）报文结构体
type NativeRequestXml struct {
	XMLName  			xml.Name 				`xml:"xml"`
	Appid 				string 					`xml:"appid"`					//公众账号ID
	Mch_id 				string 					`xml:"mch_id"`					//商户号
	Out_trade_no 		string 					`xml:"out_trade_no"`			//商户订单号
	Trade_type 			string 					`xml:"trade_type"`			//交易类型
	Fee_type 			string 					`xml:"fee_type"`				//标价币种
	Total_fee 			string 					`xml:"total_fee"`				//标价金额
	Body 				string 					`xml:"body"`					//商品描述
	Spbill_create_ip 	string 					`xml:"spbill_create_ip"`		//终端IP
	Notify_url 			string 					`xml:"notify_url"`			//通知地址
	Nonce_str 			string 					`xml:"nonce_str"`				//随机字符串
	Sign 				string 					`xml:"sign"`					//签名
}











