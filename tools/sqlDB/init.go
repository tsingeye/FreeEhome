package sqlDB

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/tsingeye/FreeEhome/tools"
	"github.com/tsingeye/FreeEhome/tools/logs"
	"os"
	"time"
)

var (
	gormDB *gorm.DB
	dbType string
)

func InitDB() {
	var err error
	var dataSource string
	//获取数据库类型
	dbType = beego.AppConfig.String("dbType")
	if dbType == "" || dbType != "mysql" && dbType != "sqlite3" {
		dbType = "sqlite3"
	}
	//获取配置文件中的数据库名
	dbName := beego.AppConfig.String(tools.StringsJoin(dbType, "::dbName"))
	switch dbType {
	case "mysql":
		dbUser := beego.AppConfig.String(tools.StringsJoin(dbType, "::dbUser"))
		dbPwd := beego.AppConfig.String(tools.StringsJoin(dbType, "::dbPwd"))
		dbAddr := beego.AppConfig.String(tools.StringsJoin(dbType, "::dbAddr"))
		dbCharset := beego.AppConfig.String(tools.StringsJoin(dbType, "::dbCharset"))

		//dataSource = tools.StringsJoin(dbUser, ":", dbPwd, "@tcp(", dbAddr, ")/", "?charset=", dbCharset, "&parseTime=True&loc=Local")
		dataSource = tools.StringsJoin(dbUser, ":", dbPwd, "@tcp(", dbAddr, ")/", "?charset=", dbCharset)
		gormDB, err = gorm.Open(dbType, dataSource)
		if err != nil {
			logs.PanicLogger.Panicln(fmt.Sprintf("failed to connect %s database: %s", dbType, err))
		}
		//判断MySQL是否创建该数据库，若无则创建
		if err = gormDB.Exec(fmt.Sprintf("CREATE DATABASE if not exists `%s`", dbName)).Error; err != nil {
			logs.PanicLogger.Panicln("failed to create database error: ", err)
		}

		dataSource = tools.StringsJoin(dbUser, ":", dbPwd, "@tcp(", dbAddr, ")/", dbName, "?charset=", dbCharset, "&parseTime=True&loc=Local")
	case "sqlite3":
		//作为服务启动时需切换到当前工作目录
		os.Chdir(tools.GetAbsPath())
		dataSource = dbName
	}

	gormDB, err = gorm.Open(dbType, dataSource)
	if err != nil {
		logs.PanicLogger.Panicln(fmt.Sprintf("failed to connect %s database: %s", dbType, err))
	}

	//最大连接周期，超过时间的连接就close
	gormDB.DB().SetConnMaxLifetime(100 * time.Second)
	//设置最大连接数，默认最大连接数为100
	gormDB.DB().SetMaxOpenConns(100)
	//设置闲置连接数
	gormDB.DB().SetMaxIdleConns(10)
	//全局禁用表名复数
	//如果设置为true,`User`的默认表名为`user`,使用`TableName`设置的表名不受影响
	gormDB.SingularTable(true)
	//启用Logger，显示详细日志
	//gormDB.LogMode(true)

	logs.BeeLogger.Info(fmt.Sprintf("successful connection to %s database", dbType))
	fmt.Printf("%s successful connection to %s database\n", time.Now().Format("2006-01-02 15:04:05"), dbType)

	//初始化数据表格
	InitDeviceList()
	InitChannelList()
}
