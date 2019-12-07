package controllers

import "encoding/xml"

//微信支付（Native支付）下单请求报文结构体
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

//微信支付（Native支付）下单返回报文结构体
type NativeResponseXml struct {
	XMLName  			xml.Name 				`xml:"xml"`
	Return_code  		string					`xml:"return_code"`			//返回状态码
	Return_msg			string					`xml:"return_msg"`			//返回信息
	Appid				string					`xml:"appid"`					//公众账号ID
	Mch_id				string					`xml:"mch_id"`					//商户号
	Device_info			string					`xml:"device_info"`			//设备号
	Nonce_str			string					`xml:"nonce_str"`				//随机字符串
	Sign				string					`xml:"sign"`					//签名
	Result_code			string					`xml:"result_code"`			//业务结果
	Err_code			string					`xml:"err_code"`				//错误代码
	Err_code_des		string					`xml:"err_code_des"`			//错误代码描述
	Trade_type			string					`xml:"trade_type"`			//交易类型
	Prepay_id			string					`xml:"prepay_id"`				//预支付交易会话标识
	Code_url 			string					`xml:"code_url"`				//二维码链接
}

//微信支付（Native支付）微信官方异步通知
type NativePayNotifyResult struct {
	XMLName				xml.Name				`xml:"xml"`
	return_code
	return_msg
	appid
	mch_id
	device_info
	nonce_str
	sign
	sign_type
	result_code
	err_code
	err_code_des
	openid
	is_subscribe
	trade_type
	bank_type
	total_fee
	settlement_total_fee
	fee_type
	cash_fee
	cash_fee_type
	transaction_id
	out_trade_no
	attach
	time_end
}









