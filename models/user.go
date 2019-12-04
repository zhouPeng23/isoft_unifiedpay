package models

type User struct {
	Id int64
	UserId string `orm:"unique"`
	UserName string
	PassWord string
	Status bool
	LoginTimes int64
}





