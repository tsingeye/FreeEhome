package config

import (
	"encoding/xml"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/tsingeye/FreeEhome/tools"
	"github.com/tsingeye/FreeEhome/tools/logs"
	"io/ioutil"
	"os"
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
	FreeEHomeCodeMap  map[int64]string
	UDPAddr           string   //udp服务地址
	STSAddr           string   //STS服务地址
	StreamStartIP     string   //实时直播流推送配置IP
	StreamStartPort   int64    //实时直播流推送配置Port
	WaitStreamURLTime int64    //等待实时直播URL超时时间
	StreamIP          []string //用于组合实时直播播放URL
	HLSPort           []string
	RTMPPort          []string
	RTSPPort          []string
	XMLConfigInfo     XMLConfig //解析配置文件中的XML文件
)

func init() {
	FreeEHomeCodeMap = map[int64]string{
		FreeEHomeSuccessOK:               "Success OK",
		FreeEHomeOperationFailed:         "Operation Failed",
		FreeEHomeParameterError:          "Parameter Error",
		FreeEHomeUnauthorized:            "Unauthorized or invalid authCode",
		FreeEHomeDeviceNotOnline:         "Device Not Online",
		FreeEHomeRequestTimeout:          "Request Timeout With STS",
		FreeEHomeServerError:             "Server Error",
		FreeEHomeChannelIDNotFound:       "ChannelID Not Found",
		FreeEHomeChannelIDNotStreamStart: "ChannelID's StreamStart Not Started",
	}

	UDPAddr = beego.AppConfig.String("udpAddr")
	if UDPAddr == "" {
		logs.PanicLogger.Panicln("init udpAddr error, udpAddr cannot be empty")
	}

	STSAddr = beego.AppConfig.String("stsAddr")
	if STSAddr == "" {
		logs.PanicLogger.Panicln("init stsAddr error, stsAddr cannot be empty")
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

	WaitStreamURLTime, err = beego.AppConfig.Int64("waitStreamURLTime")
	if err != nil || WaitStreamURLTime < 0 {
		WaitStreamURLTime = 4
		logs.BeeLogger.Error("init waitStreamURLTime, set default value waitStreamURLTime=%d", WaitStreamURLTime)
	}

	StreamIP = beego.AppConfig.Strings("streamIP")
	HLSPort = beego.AppConfig.Strings("hlsPort")
	RTMPPort = beego.AppConfig.Strings("rtmpPort")
	RTSPPort = beego.AppConfig.Strings("rtspPort")

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
