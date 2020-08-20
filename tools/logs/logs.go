package logs

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/tsingeye/FreeEhome/tools"
	"log"
	"os"
	"time"
)

var (
	//记录崩溃信息的日志，golang自带的log包
	PanicLogger *log.Logger
	//创建一个日志对象
	BeeLogger *logs.BeeLogger
)

func init() {
	appPath := tools.GetAbsPath()
	//从配置文件获取panicFile
	panicFile := beego.AppConfig.String("panicFile")
	if panicFile == "" {
		panicFile = "fELogs/panic.log"
	}

	//根据文件路径创建对应的文件夹
	panicFile = tools.MkdirAllFile(panicFile, appPath)

	//初始化panic信息日志
	initPanicLogger(panicFile)
	//初始化基于Beego的日志模块
	initBeeLogger(appPath)

	BeeLogger.Info("init logs successful!")
	fmt.Println(time.Now().Format("2006-01-02 15:04:05"), "init logs successful!!!")
}

//基于Golang官方库log来记录panic信息
func initPanicLogger(panicFile string) {
	file, err := os.OpenFile(panicFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalln("failed to OpenFile() error: ", err)
	}
	//参数二可区分日志级别，如"[Info]"或者"[Warning]"，参数三打印全路径文件名和行号
	PanicLogger = log.New(file, "", log.Llongfile)
	//设置写入文件格式：log.LstdFlags打印日期和时间
	PanicLogger.SetFlags(log.LstdFlags | log.Llongfile)
}

//基于Beego的日志库
func initBeeLogger(appPath string) {
	//创建一个日志记录器，参数为缓冲区大小
	BeeLogger = logs.NewLogger(10000)

	//获取日志输出引擎，控制台输出为console，文件输出为file
	adapterType := beego.AppConfig.String("adapterType")
	if adapterType == "" || adapterType != "file" && adapterType != "console" {
		//默认输出到file
		adapterType = "file"
	}

	//获取配置文件中的日志级别
	level, err := beego.AppConfig.Int(tools.StringsJoin(adapterType, "::level"))
	if err != nil || level < 0 || level > 7 {
		PanicLogger.Panicln("init logs's level error: ", err)
	}

	switch adapterType {
	case "console":
		//输出到控制台
		//是否开启打印日志彩色打印(需环境支持彩色打印)
		color, err := beego.AppConfig.Bool(tools.StringsJoin(adapterType, "::color"))
		if err != nil {
			PanicLogger.Panicln("init logs's color error: ", err)
		}

		config := struct {
			Level int  `json:"level"`
			Color bool `json:"color"`
		}{
			Level: level,
			Color: color,
		}
		jsonConfig, err := json.Marshal(config)
		if err != nil {
			PanicLogger.Panicln("json.Marshal() error: ", err)
		}
		//设置日志记录方式：控制台输出
		BeeLogger.SetLogger(adapterType, string(jsonConfig))
	case "file":
		//获取文件日志目录
		fileLogs := beego.AppConfig.String("fileLogs")
		if fileLogs == "" {
			fileLogs = "fELogs/fe.log"
		}
		fileLogs = tools.MkdirAllFile(fileLogs, appPath)

		//输出到文件
		//获取每个文件保存的最大尺寸
		maxSize, err := beego.AppConfig.Int64(tools.StringsJoin(adapterType, "::maxSize"))
		if err != nil || maxSize <= 0 {
			PanicLogger.Panicln("init logs's maxSize error: ", err)
		}

		//获取日志文件最多保存的天数
		maxDays, err := beego.AppConfig.Int64(tools.StringsJoin(adapterType, "::maxDays"))
		if err != nil || maxDays <= 0 {
			PanicLogger.Panicln("init logs's maxDays error: ", err)
		}

		config := struct {
			FileName string `json:"filename"`
			MaxSize  int64  `json:"maxsize"`
			MaxDays  int64  `json:"maxdays"`
		}{
			FileName: fileLogs,
			MaxSize:  maxSize * 1024 * 1024,
			MaxDays:  maxDays,
		}
		jsonConfig, err := json.Marshal(config)
		if err != nil {
			PanicLogger.Panicln("json.Marshal() error: ", err)
		}
		//设置日志记录方式：本地文件记录
		BeeLogger.SetLogger(adapterType, string(jsonConfig))
		//设置日志写入缓冲区的等级
		BeeLogger.SetLevel(level)
		////将日志从缓冲区读出，写入到文件
		BeeLogger.Flush()
	}
	//输出log时能显示输出文件名和行号
	BeeLogger.EnableFuncCallDepth(true)
}
