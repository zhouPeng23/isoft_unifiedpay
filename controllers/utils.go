package controllers

import (
	"fmt"
	"math/rand"
	"time"
	"sync"
)

var (
	intRandom = 100000000
	lock sync.Mutex
)


//生成订单随机数
func QueryUniqueRandom() string {
	if intRandom == 999999999{
		intRandom = 100000000
	}
	lock.Lock()
	intRandom++
	lock.Unlock()
	//获取一个长度为9的唯一数字字符串（给自己人看的）
	intStr9 := "000000000" + fmt.Sprintf("%d",intRandom)
	intStr9 = intStr9[len(intStr9)-9:len(intStr9)]
	//获取一个长度为7的随机数（用于干扰别有用心者）
	rand.Seed(time.Now().UnixNano())
	random := rand.Intn(7777777)
	strRandom7 := "0000000" + fmt.Sprintf("%d", random)
	strRandom7 = strRandom7[len(strRandom7)-7:len(strRandom7)]
	return intStr9 + strRandom7
}



