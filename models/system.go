package models

import (
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

//登录
func Login(username, password, remoteAddr string) (replyData map[string]interface{}) {
	//授权码使用用户名+密码+UUID使用MD5加密生成
	authCode := tools.GetMD5String(tools.StringsJoin(username, password, tools.GetUUID()))
	authClient := &config.AuthClient{
		Username:   username,
		Password:   password,
		RemoteAddr: remoteAddr,
		AuthCode:   authCode,
	}
	//登录成功，记录authCode
	config.AuthCheck.Set(authCode, authClient, cache.DefaultExpiration)

	logs.BeeLogger.Info("username=%s, password=%s login successfully, get authCode=%s", username, password, authCode)

	replyData = map[string]interface{}{
		"errCode":  config.FreeEHomeSuccessOK,
		"errMsg":   config.FreeEHomeCodeMap[config.FreeEHomeSuccessOK],
		"authCode": authCode,
	}
	return
}

//登出
func Logout(authCode string) (replyData map[string]interface{}) {
	logs.BeeLogger.Info("user logout successfully, delete authCode=%s", authCode)
	//登出成功，删除authCode
	config.AuthCheck.Delete(authCode)

	replyData = map[string]interface{}{
		"errCode": config.FreeEHomeSuccessOK,
		"errMsg":  config.FreeEHomeCodeMap[config.FreeEHomeSuccessOK],
	}
	return
}

//获取系统信息：CPU、网络、内存等v
func SystemInfo(authCode string) (replyData map[string]interface{}) {
	replyData = map[string]interface{}{
		"errCode":  config.FreeEHomeServerError,
		"errMsg":   config.FreeEHomeCodeMap[config.FreeEHomeServerError],
		"authCode": authCode,
	}

	//CPU信息
	percent, err := cpu.Percent(time.Second, false)
	if err != nil {
		logs.BeeLogger.Error("authCode=%s, cpu.Percent error: %s", authCode, err)
		return
	}

	//内存信息
	memory, err := mem.VirtualMemory()
	if err != nil {
		logs.BeeLogger.Error("authCode=%s, mem.VirtualMemory error: %s", authCode, err)
		return
	}

	//网络信息
	beforeIO, err := net.IOCounters(false)
	if err != nil {
		logs.BeeLogger.Error("authCode=%s, net.IOCounters error: %s", authCode, err)
		return
	}
	time.Sleep(1 * time.Second)
	afterIO, err := net.IOCounters(false)
	if err != nil {
		logs.BeeLogger.Error("authCode=%s, net.IOCounters error: %s", authCode, err)
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
		"authCode":       authCode,
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
