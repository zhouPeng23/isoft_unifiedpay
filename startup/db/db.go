package db

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql" // _ 的作用,并不需要把整个包都导入进来,仅仅是是希望它执行init()函数而已
	"net/url"
	"isoft_unifiedpay/models"
	"isoft_unifiedpay/common/chiperutil"
	"github.com/astaxie/beego/logs"
)

func init() {
	dbhost := beego.AppConfig.String("db.host")
	dbport := beego.AppConfig.String("db.port")
	dbname := beego.AppConfig.String("db.name")
	dbuser := beego.AppConfig.String("db.user")
	dbpass := beego.AppConfig.String("db.pass")
	timezone := beego.AppConfig.String("db.timezone")
	aesChiperKey := beego.AppConfig.String("isoft_unifiedpay.aes.cipher.key")
	// 对数据库密码进行解密
	dbport = chiperutil.AesDecryptToStr(dbport, aesChiperKey)
	dbuser = chiperutil.AesDecryptToStr(dbuser, aesChiperKey)
	dbpass = chiperutil.AesDecryptToStr(dbpass, aesChiperKey)

	Dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?allowNativePasswords=true&charset=utf8", dbuser, dbpass, dbhost, dbport, dbname)

	if timezone != "" {
		Dsn = Dsn + "&loc=" + url.QueryEscape(timezone)
	}

	orm.RegisterDriver("mysql", orm.DRMySQL)
	orm.RegisterDataBase("default", "mysql", Dsn)
	orm.SetMaxIdleConns("default", 100) // SetMaxIdleConns用于设置闲置的连接数
	orm.SetMaxOpenConns("default", 200) // SetMaxOpenConns用于设置最大打开的连接数,默认值为0表示不限制
	db, _ := orm.GetDB("default")
	db.SetConnMaxLifetime(100)

	if beego.AppConfig.String("runmode") == "dev" {
		orm.Debug = true
	}
	registerModel()
	createTable()
}

func registerModel() {
	orm.RegisterModel(new(models.Order))
	orm.RegisterModel(new(models.User))
}

// 自动建表
func createTable() {
	force := false  // 不强制建数据库
	verbose := true // 打印建表过程
	if err := orm.RunSyncdb("default", force, verbose); err != nil {
		logs.Error(err)
	}
}
