package models

import (
	"fmt"
	"github.com/tsingeye/FreeEhome/config"
	"github.com/tsingeye/FreeEhome/service/udp"
	"github.com/tsingeye/FreeEhome/tools"
	"github.com/tsingeye/FreeEhome/tools/logs"
	"github.com/tsingeye/FreeEhome/tools/sqlDB"
	"time"
)

//先判断内存中是否存在可用的sessionURL，若存在则直接返回
func StartStream(token, channelID string) (replyData map[string]interface{}) {
	sessionInfo := config.StreamSession.Get(channelID)
	if sessionInfo != nil && sessionInfo.SessionURL != nil {
		replyData = map[string]interface{}{
			"errCode":    config.FreeEHomeSuccessOK,
			"errMsg":     config.FreeEHomeCodeMap[config.FreeEHomeSuccessOK],
			"sessionURL": sessionInfo.SessionURL,
		}
		return
	}

	var channel sqlDB.ChannelList
	if !sqlDB.First(&channel, "ChannelID = ?", channelID) {
		replyData = map[string]interface{}{
			"errCode": config.FreeEHomeParameterError,
			"errMsg":  config.FreeEHomeCodeMap[config.FreeEHomeParameterError],
		}
		return
	}

	errCode, session := udp.SendInviteStream(token, channel.DeviceID, channelID)

	replyData = map[string]interface{}{
		"errCode": errCode,
		"errMsg":  config.FreeEHomeCodeMap[errCode],
	}

	if errCode == config.FreeEHomeSuccessOK {
		fmt.Printf("%s token=%s, DeviceID=%s, ChannelID=%s, get session=%s\n", time.Now().Format("2006-01-02 15:04:05"), token, channel.DeviceID, channelID, session)
		logs.BeeLogger.Emergency("token=%s, DeviceID=%s, ChannelID=%s, get session=%s", token, channel.DeviceID, channelID, session)

		hexSession, ok := config.HookSession.Get(session)
		if ok {
			//说明已经收到hook返回的对应的session，此时无需再等待，可直接合成sessionURL
			//删除hook返回的session
			config.HookSession.Delete(session)

			sessionURL := getSessionURL(hexSession.(string))
			sessionInfo = &config.SessionInfo{
				Session:    session,
				DeviceID:   channel.DeviceID,
				ChannelID:  channelID,
				SessionURL: sessionURL,
			}
			config.StreamSession.Set(channelID, sessionInfo)
			replyData["sessionURL"] = sessionURL
		} else {
			retBool, hexSession := waitSessionFromHook(session)
			if retBool {
				//在超时时间内等到hook返回的session，此时可直接合成sessionURL
				//删除hook返回的session
				config.HookSession.Delete(session)
				sessionURL := getSessionURL(hexSession)
				sessionInfo = &config.SessionInfo{
					Session:    session,
					DeviceID:   channel.DeviceID,
					ChannelID:  channelID,
					SessionURL: sessionURL,
				}
				config.StreamSession.Set(channelID, sessionInfo)
				replyData["sessionURL"] = sessionURL
			} else {
				//在规定的时间内未等到hook返回的session，需清除内存中的session
				go udp.SendByeStream(token, channel.DeviceID, channelID)

				replyData = map[string]interface{}{
					"errCode": config.FreeEHomeRequestTimeout,
					"errMsg":  config.FreeEHomeCodeMap[config.FreeEHomeRequestTimeout],
				}
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
		retMap["rtmp"] = tools.StringsJoin("rtmp://", config.StreamIP[0], ":", config.RTMPPort[0], "/rtp/", session)
		retMap["flv"] = tools.StringsJoin("http://", config.StreamIP[0], ":", config.HLSPort[0], "/rtp/", session, ".flv")
		retMap["rtsp"] = tools.StringsJoin("rtsp://", config.StreamIP[0], ":", config.RTSPPort[0], "/rtp/", session)
		retMap["hls"] = tools.StringsJoin("http://", config.StreamIP[0], ":", config.HLSPort[0], "/rtp/", session, "/hls.m3u8")
	default:
		//配置文件中streamIP有多个IP地址
		retMap["rtmp"] = tools.StringsJoin("rtmp://", config.StreamIP[1], ":", config.RTMPPort[1], "/rtp/", session)
		retMap["flv"] = tools.StringsJoin("http://", config.StreamIP[1], ":", config.HLSPort[1], "/rtp/", session, ".flv")
		retMap["rtsp"] = tools.StringsJoin("rtsp://", config.StreamIP[1], ":", config.RTSPPort[1], "/rtp/", session)
		retMap["hls"] = tools.StringsJoin("http://", config.StreamIP[1], ":", config.HLSPort[1], "/rtp/", session, "/hls.m3u8")

	}

	return retMap
}

//获取到设备返回的session后等待hook返回的对应的session，若在超时时间内返回则成功
func waitSessionFromHook(session string) (bool, string) {
	idleDelay := time.NewTimer(time.Duration(config.WaitHookSessionTime) * time.Second)
	defer idleDelay.Stop()

	for {
		select {
		case <-idleDelay.C:
			//超时
			return false, ""
		default:
			//每隔100毫秒请求一次数据
			time.Sleep(100 * time.Millisecond)

			hexSession, ok := config.HookSession.Get(session)
			if ok {
				return true, hexSession.(string)
			}
		}
	}
}
