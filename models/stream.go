package models

import (
	"fmt"
	"github.com/tsingeye/FreeEhome/config"
	"github.com/tsingeye/FreeEhome/service/udp"
	"github.com/tsingeye/FreeEhome/tools"
	"github.com/tsingeye/FreeEhome/tools/logs"
	"time"
)

//先判断内存中是否存在可用的sessionURL，若存在则直接返回
func StartStream(authCode, deviceID, channelID string) (replyData map[string]interface{}) {
	sessionInfo := config.StreamSession.Get(channelID)
	if sessionInfo != nil && sessionInfo.SessionURL != nil {
		replyData = map[string]interface{}{
			"errCode":    config.FreeEHomeSuccessOK,
			"errMsg":     config.FreeEHomeCodeMap[config.FreeEHomeSuccessOK],
			"authCode":   authCode,
			"sessionURL": sessionInfo.SessionURL,
		}
		return
	}

	errCode, session := udp.SendInviteStream(authCode, deviceID, channelID)

	replyData = map[string]interface{}{
		"errCode":  errCode,
		"errMsg":   config.FreeEHomeCodeMap[errCode],
		"authCode": authCode,
	}

	if errCode == config.FreeEHomeSuccessOK {
		fmt.Printf("%s authCode=%s, DeviceID=%s, ChannelID=%s, get session=%s\n", time.Now().Format("2006-01-02 15:04:05"), authCode, deviceID, channelID, session)
		logs.BeeLogger.Emergency("authCode=%s, DeviceID=%s, ChannelID=%s, get session=%s", authCode, deviceID, channelID, session)

		_, ok := config.HookSession.Get(session)
		if ok {
			//说明已经收到hook返回的对应的session，此时无需再等待，可直接合成sessionURL
			//删除hook返回的session
			config.HookSession.Delete(session)
			sessionURL := getSessionURL(session)
			sessionInfo = &config.SessionInfo{
				Session:    session,
				DeviceID:   deviceID,
				ChannelID:  channelID,
				SessionURL: sessionURL,
			}
			config.StreamSession.Set(channelID, sessionInfo)
			replyData["sessionURL"] = sessionURL
		} else if waitSessionFromHook(session) {
			//在超时时间内等到hook返回的session，此时可直接合成sessionURL
			//删除hook返回的session
			config.HookSession.Delete(session)
			sessionURL := getSessionURL(session)
			sessionInfo = &config.SessionInfo{
				Session:    session,
				DeviceID:   deviceID,
				ChannelID:  channelID,
				SessionURL: sessionURL,
			}
			config.StreamSession.Set(channelID, sessionInfo)
			replyData["sessionURL"] = sessionURL
		} else {
			//在规定的时间内未等到hook返回的session，需清除内存中的session
			go udp.SendByeStream(authCode, deviceID, channelID)

			replyData = map[string]interface{}{
				"errCode":  config.FreeEHomeRequestTimeout,
				"errMsg":   config.FreeEHomeCodeMap[config.FreeEHomeRequestTimeout],
				"authCode": authCode,
			}
		}
	}
	return
}

//生成返回的sessionURL
func getSessionURL(session string) (retMap map[string]string) {
	retMap = make(map[string]string)
	switch len(config.StreamIP) {
	case 1:
		//配置文件中streamIP只有一个IP地址
		retMap["rtmp"] = tools.StringsJoin("rtmp://", config.StreamIP[0], ":", config.RTMPPort[0], "/live/", session)
		retMap["flv"] = tools.StringsJoin("http://", config.StreamIP[0], ":", config.HLSPort[0], "/live/", session, ".flv")
		retMap["rtsp"] = tools.StringsJoin("rtsp://", config.StreamIP[0], ":", config.RTSPPort[0], "/live/", session)
		retMap["hls"] = tools.StringsJoin("http://", config.StreamIP[0], ":", config.HLSPort[0], "/live/", session, "/hls.m3u8")
	default:
		//配置文件中streamIP有多个IP地址
		retMap["rtmp"] = tools.StringsJoin("rtmp://", config.StreamIP[1], ":", config.RTMPPort[1], "/live/", session)
		retMap["flv"] = tools.StringsJoin("http://", config.StreamIP[1], ":", config.HLSPort[1], "/live/", session, ".flv")
		retMap["rtsp"] = tools.StringsJoin("rtsp://", config.StreamIP[1], ":", config.RTSPPort[1], "/live/", session)
		retMap["hls"] = tools.StringsJoin("http://", config.StreamIP[1], ":", config.HLSPort[1], "/live/", session, "/hls.m3u8")

	}

	return nil
}

//获取到设备返回的session后等待hook返回的对应的session，若在超时时间内返回则成功
func waitSessionFromHook(session string) bool {
	idleDelay := time.NewTimer(time.Duration(config.WaitHookSessionTime) * time.Second)
	defer idleDelay.Stop()

	for {
		select {
		case <-idleDelay.C:
			//超时
			return false
		default:
			//每隔100毫秒请求一次数据
			time.Sleep(100 * time.Millisecond)

			_, ok := config.HookSession.Get(session)
			if ok {
				return true
			}
		}
	}
}

func StopStream(authCode, deviceID, channelID string) (replyData map[string]interface{}) {
	replyData = map[string]interface{}{
		"errCode":  config.FreeEHomeSuccessOK,
		"errMsg":   config.FreeEHomeCodeMap[config.FreeEHomeSuccessOK],
		"authCode": authCode,
	}
	return
}
