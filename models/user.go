package models

type User struct {
	Id int64
	UserId 						string				`orm:"unique"`					//用户ID
	UserName 					string												//用户名
	PassWord 					string												//密码
	Status 						bool												//用户状态
	LoginIp						string												//最近一次登录ip
	LoginTimes 					int64												//登录次数
}





