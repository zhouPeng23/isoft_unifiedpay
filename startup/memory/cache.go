package memory

import "github.com/astaxie/beego/cache"

var CacheEngine cache.Cache

func init()  {
	CacheEngine, _  = cache.NewCache("memory", `{"interval":60}`)
}


