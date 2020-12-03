package models

import (
	"encoding/json"
	"fmt"
	"github.com/patrickmn/go-cache"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
	"github.com/tsingeye/FreeEhome/config"
	"github.com/tsingeye/FreeEhome/tools"
	"github.com/tsingeye/FreeEhome/tools/logs"
	"github.com/tsingeye/FreeEhome/tools/sqlDB"
	"time"
)

//用于解析用户登录数据
type userLogin struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

//登录
func Login(remoteAddr string, requestBody []byte) (replyData map[string]interface{}) {
	var loginData userLogin
	if err := json.Unmarshal(requestBody, &loginData); err != nil {
		logs.BeeLogger.Error("SystemController.Login() ---> remoteAddr=%s, json.Unmarshal() error:%s", remoteAddr, err)
		replyData = map[string]interface{}{
			"errCode": config.FreeEHomeParameterError,
			"errMsg":  config.FreeEHomeCodeMap[config.FreeEHomeParameterError],
		}

		return
	}

	//授权码使用用户名+密码+UUID使用MD5加密生成
	token := tools.GetMD5String(tools.StringsJoin(loginData.UserName, loginData.Password, tools.GetUUID()))

	//登录成功，记录authCode
	config.AuthCheck.Set(token, token, cache.DefaultExpiration)

	logs.BeeLogger.Info("username=%s, password=%s login successfully, get token=%s", loginData.UserName, loginData.Password, token)

	replyData = map[string]interface{}{
		"errCode": config.FreeEHomeSuccessOK,
		"errMsg":  config.FreeEHomeCodeMap[config.FreeEHomeSuccessOK],
		"token":   token,
	}
	return
}

//登出
func Logout(token string) (replyData map[string]interface{}) {
	logs.BeeLogger.Info("user logout successfully, delete token=%s", token)
	//登出成功，删除token
	config.AuthCheck.Delete(token)

	replyData = map[string]interface{}{
		"errCode": config.FreeEHomeSuccessOK,
		"errMsg":  config.FreeEHomeCodeMap[config.FreeEHomeSuccessOK],
	}
	return
}

//获取系统信息：CPU、网络、内存等v
func SystemInfo() (replyData map[string]interface{}) {
	replyData = map[string]interface{}{
		"errCode":  config.FreeEHomeServerError,
		"errMsg":   config.FreeEHomeCodeMap[config.FreeEHomeServerError],
	}

	//CPU信息
	percent, err := cpu.Percent(time.Second, false)
	if err != nil {
		logs.BeeLogger.Error("cpu.Percent error: %s", err)
		return
	}

	//内存信息
	memory, err := mem.VirtualMemory()
	if err != nil {
		logs.BeeLogger.Error("mem.VirtualMemory error: %s", err)
		return
	}

	//网络信息
	beforeIO, err := net.IOCounters(false)
	if err != nil {
		logs.BeeLogger.Error("net.IOCounters error: %s", err)
		return
	}
	time.Sleep(1 * time.Second)
	afterIO, err := net.IOCounters(false)
	if err != nil {
		logs.BeeLogger.Error("net.IOCounters error: %s", err)
		return
	}

	//上传速度
	uploadSpeed := (afterIO[0].BytesSent - beforeIO[0].BytesSent) / 1024
	//fmt.Println("打印上传速度：", afterIO[0].BytesSent - beforeIO[0].BytesSent)
	//下载速度
	downloadSpeed := (afterIO[0].BytesRecv - beforeIO[0].BytesRecv) / 1024
	//fmt.Println("打印下载速度：", afterIO[0].BytesRecv - beforeIO[0].BytesRecv)

	replyData = map[string]interface{}{
		"errCode":        config.FreeEHomeSuccessOK,
		"errMsg":         config.FreeEHomeCodeMap[config.FreeEHomeSuccessOK],
		"cpuUsedPercent": fmt.Sprintf("%d%s", int(percent[0]), "%"),
		"virtualMemory": map[string]interface{}{
			"total":       fmt.Sprintf("%dMB", memory.Total/1024/1024),
			"available":   fmt.Sprintf("%dMB", memory.Available/1024/1024),
			"used":        fmt.Sprintf("%dMB", memory.Used/1024/1024),
			"usedPercent": fmt.Sprintf("%d%s", int(memory.UsedPercent), "%"),
		},
		"network": map[string]interface{}{
			"uploadSpeed":   fmt.Sprintf("%dKB/s", uploadSpeed),
			"downloadSpeed": fmt.Sprintf("%dKB/s", downloadSpeed),
		},
		"deviceInfo": map[string]interface{}{
			"totalCount":  sqlDB.Count(&sqlDB.DeviceList{}, ""),
			"onlineCount": sqlDB.Count(&sqlDB.DeviceList{}, "Status = ?", "ON"),
		},
		"channelInfo": map[string]interface{}{
			"totalCount":  sqlDB.Count(&sqlDB.ChannelList{}, ""),
			"onlineCount": sqlDB.Count(&sqlDB.ChannelList{}, "Status = ?", "ON"),
		},
	}
	return
}
