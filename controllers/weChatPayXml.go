package controllers

import "encoding/xml"

//微信支付（Native支付）下单请求报文结构体
type NativeRequestXml struct {
	XMLName  					xml.Name 				`xml:"xml"`
	Appid 						string 					`xml:"appid"`					//公众账号ID
	Mch_id 						string 					`xml:"mch_id"`					//商户号
	Out_trade_no 				string 					`xml:"out_trade_no"`			//商户订单号
	Trade_type 					string 					`xml:"trade_type"`			//交易类型
	Fee_type 					string 					`xml:"fee_type"`				//标价币种
	Total_fee 					string 					`xml:"total_fee"`				//标价金额
	Body 						string 					`xml:"body"`					//商品描述
	Spbill_create_ip 			string 					`xml:"spbill_create_ip"`		//终端IP
	Notify_url 					string 					`xml:"notify_url"`			//通知地址
	Nonce_str 					string 					`xml:"nonce_str"`				//随机字符串
	Sign 						string 					`xml:"sign"`					//签名
}

//微信支付（Native支付）下单返回报文结构体
type NativeResponseXml struct {
	XMLName  					xml.Name 				`xml:"xml"`
	Return_code  				string					`xml:"return_code"`			//返回状态码
	Return_msg					string					`xml:"return_msg"`			//返回信息
	Appid						string					`xml:"appid"`					//公众账号ID
	Mch_id						string					`xml:"mch_id"`					//商户号
	Device_info					string					`xml:"device_info"`			//设备号
	Nonce_str					string					`xml:"nonce_str"`				//随机字符串
	Sign						string					`xml:"sign"`					//签名
	Result_code					string					`xml:"result_code"`			//业务结果
	Err_code					string					`xml:"err_code"`				//错误代码
	Err_code_des				string					`xml:"err_code_des"`			//错误代码描述
	Trade_type					string					`xml:"trade_type"`			//交易类型
	Prepay_id					string					`xml:"prepay_id"`				//预支付交易会话标识
	Code_url 					string					`xml:"code_url"`				//二维码链接
}

//微信支付（Native支付）微信官方异步通知
type NativePayNotifyResult struct {
	XMLName						xml.Name				`xml:"xml"`
	Return_code					string					`xml:"return_code"`			//返回状态码
	Return_msg					string					`xml:"return_msg"`			//返回信息
	Appid						string					`xml:"appid"`					//公众账号ID
	Mch_id						string					`xml:"mch_id"`					//商户号
	Device_info					string					`xml:"device_info"`			//设备号
	Nonce_str					string					`xml:"nonce_str"`				//随机字符串
	Sign						string					`xml:"sign"`					//签名
	Sign_type					string					`xml:"sign_type"`				//签名类型
	Result_code					string					`xml:"result_code"`			//业务结果
	Err_code					string					`xml:"err_code"`				//错误代码
	Err_code_des				string					`xml:"err_code_des"`			//错误代码描述
	Openid						string					`xml:"openid"`					//用户标识
	Is_subscribe				string					`xml:"is_subscribe"`			//是否关注公众账号
	Trade_type					string					`xml:"trade_type"`			//交易类型
	Bank_type					string					`xml:"bank_type"`				//付款银行
	Total_fee					string					`xml:"total_fee"`				//订单金额
	Settlement_total_fee		string					`xml:"settlement_total_fee"`	//应结订单金额
	Fee_type					string					`xml:"fee_type"`				//货币种类
	Cash_fee					string					`xml:"cash_fee"`				//现金支付金额
	Cash_fee_type				string					`xml:"cash_fee_type"`			//现金支付货币类型
	Transaction_id				string					`xml:"transaction_id"`		//微信支付订单号
	Out_trade_no				string					`xml:"out_trade_no"`			//商户订单号
	Attach						string					`xml:"attach"`					//商家数据包
	Time_end					string					`xml:"time_end"`				//支付完成时间
}

//微信支付（Native支付）异步通知-回复
type NativePayNotifyResponse struct {
	XMLName 					xml.Name
	Return_code					string					`xml:"return_code"`			//返回状态码
	Return_msg					string					`xml:"return_msg"`			//返回信息
}







