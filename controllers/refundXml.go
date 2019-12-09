package controllers

import "encoding/xml"

//申请退款请求报文结构体
type RefundRequestXml struct {
	XMLName  					xml.Name 				`xml:"xml"`
	Appid 						string 					`xml:"appid"`						//公众账号ID
	Mch_id 						string 					`xml:"mch_id"`						//商户号
	Nonce_str					string					`xml:"nonce_str"`					//随机字符串
	Out_refund_no				string					`xml:"out_refund_no"`				//退款订单
	Refund_fee					string					`xml:"refund_fee"`				//本次退款金额
	Refund_fee_type				string					`xml:"refund_fee_type"`			//退款货币种类
	Out_trade_no				string					`xml:"out_trade_no"`				//商户订单号(原交易订单)
	Total_fee					string					`xml:"total_fee"`					//原交易订单金额
	Refund_desc					string					`xml:"refund_desc"`				//退款原因描述
	Notify_url					string					`xml:"notify_url"`				//退款结果通知url
	Sign						string					`xml:"sign"`						//签名
}

//申请退款返回报文结构体
type RefundResponseXml struct {
	XMLName  					xml.Name 				`xml:"xml"`
	Return_code  				string 					`xml:"return_code"`				//返回状态码
	Return_msg  				string 					`xml:"return_msg"`				//返回信息
	Result_code  				string 					`xml:"result_code"`				//业务结果
	Err_code	  				string 					`xml:"err_code"`					//错误代码
	Err_code_des	  			string 					`xml:"err_code_des"`				//错误代码描述
	Appid			  			string 					`xml:"appid"`						//公众账号ID
	Mch_id			  			string 					`xml:"mch_id"`						//商户号
	Nonce_str			  		string 					`xml:"nonce_str"`					//随机字符串
	Sign				  		string 					`xml:"sign"`						//签名
	Transaction_id				string 					`xml:"transaction_id"`			//微信订单号
	Out_trade_no				string 					`xml:"out_trade_no"`				//商户订单号
	Out_refund_no				string 					`xml:"out_refund_no"`				//商户退款单号
	Refund_id					string 					`xml:"refund_id"`					//微信退款单号
	Refund_fee					string 					`xml:"refund_fee"`				//退款金额
	Settlement_refund_fee		string 					`xml:"settlement_refund_fee"`	//应结退款金额
	Total_fee					string 					`xml:"total_fee"`					//标价金额
	Settlement_total_fee 		string 					`xml:"settlement_total_fee "`	//应结订单金额
	Fee_type 					string 					`xml:"fee_type "`					//标价币种
	Cash_fee 					string 					`xml:"cash_fee "`					//现金支付金额
	Cash_fee_type 				string 					`xml:"cash_fee_type "`			//现金支付币种
	Cash_refund_fee 			string 					`xml:"cash_refund_fee "`			//现金退款金额
}

//退款结果异步通知
type RefundNotifyResult struct {
	XMLName						xml.Name				`xml:"xml"`
	Return_code					string					`xml:"return_code"`				//返回状态码
	Return_msg					string					`xml:"return_msg"`				//返回信息
	Appid						string					`xml:"appid"`						//公众账号ID
	Mch_id						string					`xml:"mch_id"`						//退款的商户号
	Nonce_str					string					`xml:"nonce_str"`					//随机字符串
	Req_info					string					`xml:"req_info"`					//加密信息
}

//退款结果加密信息
type RefInfo struct {
	XMLName 					xml.Name				`xml:"root"`
	Transaction_id				string					`xml:"transaction_id"`			//微信订单号
	Out_trade_no				string					`xml:"out_trade_no"`				//商户订单号
	Refund_id					string					`xml:"refund_id"` 					//微信退款单号
	Out_refund_no				string					`xml:"out_refund_no"`				//商户退款单号
	Total_fee					string					`xml:"total_fee"`					//订单金额
	Settlement_total_fee		string					`xml:"settlement_total_fee"`		//应结订单金额
	Refund_fee					string					`xml:"refund_fee"`				//申请退款金额
	Settlement_refund_fee		string					`xml:"settlement_refund_fee"`	//退款金额
	Refund_status				string					`xml:"refund_status"`				//退款状态
	Success_time				string					`xml:"success_time"`				//退款成功时间
	Refund_recv_accout			string					`xml:"refund_recv_accout"`		//退款入账账户
	Refund_account				string					`xml:"refund_account"`			//退款资金来源
	Refund_request_source		string					`xml:"refund_request_source"`	//退款发起来源
}

//退款异步通知-回复
type RefundNotifyResponse struct {
	XMLName 					xml.Name
	Return_code					string					`xml:"return_code"`				//返回状态码
	Return_msg					string					`xml:"return_msg"`				//返回信息
}







