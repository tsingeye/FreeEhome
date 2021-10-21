package config

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/astaxie/beego"
	"github.com/tsingeye/FreeEhome/tools"
	"github.com/tsingeye/FreeEhome/tools/logs"
)

const (
	TimeLayout = "2006-01-02 15:04:05" //时间格式统一模板
)

const (
	FreeEHomeSuccessOK               = 200
	FreeEHomeOperationFailed         = 201 //操作失败
	FreeEHomeParameterError          = 300 //参数错误
	FreeEHomeUnauthorized            = 401 //未授权或者授权码失效
	FreeEHomeDeviceNotOnline         = 406 //设备离线
	FreeEHomeRequestTimeout          = 408 //请求超时
	FreeEHomeServerError             = 500 //服务器错误
	FreeEHomeChannelIDNotFound       = 600 //通道ID不存在
	FreeEHomeChannelIDNotStreamStart = 601 //该通道ID实时直播未启动
)

var (
	FreeEHomeCodeMap      map[int64]string
	UDPAddr               string   //udp服务地址
	AlarmAddr             string   //alarm服务地址
	StreamStartIP         string   //实时直播流推送配置IP
	StreamStartPort       int64    //实时直播流推送配置Port
	WaitStreamSessionTime int64    //等待启动实时直播设备返回session的超时时间
	WaitHookSessionTime   int64    //等待hook返回对应的session的超时时间
	StreamIP              []string //用于组合实时直播播放URL
	HLSPort               []string
	RTMPPort              []string
	RTSPPort              []string
	XMLConfigInfo         XMLConfig //解析配置文件中的XML文件
)

func init() {
	FreeEHomeCodeMap = map[int64]string{
		FreeEHomeSuccessOK:               "Success OK",
		FreeEHomeOperationFailed:         "Operation Failed",
		FreeEHomeParameterError:          "Parameter Error",
		FreeEHomeUnauthorized:            "Unauthorized or invalid authCode",
		FreeEHomeDeviceNotOnline:         "Device Not Online",
		FreeEHomeRequestTimeout:          "Request Timeout",
		FreeEHomeServerError:             "Server Error",
		FreeEHomeChannelIDNotFound:       "ChannelID Not Found",
		FreeEHomeChannelIDNotStreamStart: "ChannelID's StreamStart Not Started",
	}

	UDPAddr = beego.AppConfig.String("udpAddr")
	if UDPAddr == "" {
		logs.PanicLogger.Panicln("init udpAddr error, udpAddr cannot be empty")
	}

	AlarmAddr = beego.AppConfig.String("alarmAddr")
	if AlarmAddr == "" {
		logs.PanicLogger.Panicln("init AlarmAddr error, please check it!!!")
	}

	StreamStartIP = beego.AppConfig.String("streamStartIP")
	if StreamStartIP == "" {
		logs.PanicLogger.Panicln("init streamStartIP error, streamStartIP cannot be empty")
	}
	var err error
	StreamStartPort, err = beego.AppConfig.Int64("streamStartPort")
	if err != nil {
		logs.PanicLogger.Panicln(fmt.Sprintf("init streamStartPort error: %s", err))
	}

	WaitStreamSessionTime, err = beego.AppConfig.Int64("waitStreamSessionTime")
	if err != nil || WaitStreamSessionTime < 0 {
		WaitStreamSessionTime = 3
		logs.BeeLogger.Error("init waitStreamSessionTime error, set default value waitStreamSessionTime=%d", WaitStreamSessionTime)
	}

	WaitHookSessionTime, err = beego.AppConfig.Int64("waitHookSessionTime")
	if err != nil || WaitHookSessionTime < 0 {
		WaitHookSessionTime = 3
		logs.BeeLogger.Error("init waitHookSessionTime error, set default value waitHookSessionTime=%d", WaitHookSessionTime)
	}

	StreamIP = beego.AppConfig.Strings("streamIP")
	switch len(StreamIP) {
	case 0:
		logs.PanicLogger.Panicln("init streamIP error, please check it")
	default:
		if StreamIP[0] == "" {
			logs.PanicLogger.Panicln("init streamIP error, please check it")
		}
		if len(StreamIP) >= 2 && StreamIP[1] == "" {
			logs.PanicLogger.Panicln("init streamIP error, please check it")
		}
	}
	HLSPort = beego.AppConfig.Strings("hlsPort")
	switch len(HLSPort) {
	case 0:
		logs.PanicLogger.Panicln("init hlsPort error, please check it")
	default:
		if HLSPort[0] == "" {
			logs.PanicLogger.Panicln("init hlsPort error, please check it")
		}
		if len(HLSPort) >= 2 && HLSPort[1] == "" {
			logs.PanicLogger.Panicln("init hlsPort error, please check it")
		}
	}
	RTMPPort = beego.AppConfig.Strings("rtmpPort")
	switch len(RTMPPort) {
	case 0:
		logs.PanicLogger.Panicln("init rtmpPort error, please check it")
	default:
		if RTMPPort[0] == "" {
			logs.PanicLogger.Panicln("init rtmpPort error, please check it")
		}
		if len(RTMPPort) >= 2 && RTMPPort[1] == "" {
			logs.PanicLogger.Panicln("init rtmpPort error, please check it")
		}
	}
	RTSPPort = beego.AppConfig.Strings("rtspPort")
	switch len(RTSPPort) {
	case 0:
		logs.PanicLogger.Panicln("init rtspPort error, please check it")
	default:
		if RTSPPort[0] == "" {
			logs.PanicLogger.Panicln("init rtspPort error, please check it")
		}
		if len(RTSPPort) >= 2 && RTSPPort[1] == "" {
			logs.PanicLogger.Panicln("init rtspPort error, please check it")
		}
	}

	if len(StreamIP) > 1 {
		//当区分内网和外网时，端口配置必须保证有前两个是非空值
		if len(HLSPort) <= 1 || len(RTMPPort) <= 1 || len(RTSPPort) <= 1 {
			logs.PanicLogger.Panicln("distinguish between intranet and extranet, please check it")
		}
	}

	appPath := tools.GetAbsPath()
	os.Chdir(appPath)
	file, err := os.Open("./conf/config.xml")
	if err != nil {
		logs.PanicLogger.Panicln(fmt.Sprintf("Open config.xml error, init failed, error: %s", err))
	}

	defer file.Close()
	data, err := ioutil.ReadAll(file)
	if err != nil {
		logs.PanicLogger.Panicln(fmt.Sprintf("ioutil.ReadAll() config.xml error: %s", err))
	}

	err = xml.Unmarshal(data, &XMLConfigInfo)
	if err != nil {
		logs.PanicLogger.Panicln(fmt.Sprintf("xml.Unmarshal() config.xml error: %s", err))
	}
	logs.BeeLogger.Info("init config.xml success!")
}
