package udp

import (
	"encoding/xml"
	"fmt"
	"github.com/tsingeye/FreeEhome/config"
	"github.com/tsingeye/FreeEhome/tools/logs"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

//使用DeviceID筛选出指定的UDPClient
func filterUDPClientFromDeviceID(deviceID string) *ClientForUDP {
	gw.RLock()
	udpClient, ok := gw.UDPClientList[deviceID]
	gw.RUnlock()
	if !ok {
		return nil
	}

	return udpClient
}

//发送启动实时直播流
func SendInviteStream(authCode, deviceID, channelID string) (int64, string) {
	fmt.Printf("%s authCode=%s, DeviceID=%s, ChannelID=%s send inviteStream\n", time.Now().Format("2006-01-02 15:04:05"), authCode, deviceID, channelID)
	udpClient := filterUDPClientFromDeviceID(deviceID)
	if udpClient == nil {
		return config.FreeEHomeDeviceNotOnline, ""
	}

	//fmt.Printf("%s deviceID=%s, channelID=%s, please send inviteStream in order to get a new session\n", time.Now().Format("2006-01-02 15:04:05"), deviceID, channelID)
	logs.BeeLogger.Emergency("deviceID=%s, channelID=%s, please send inviteStream in order to get a new session", deviceID, channelID)
	temp := strings.Split(channelID, "_")
	if len(temp) != 2 {
		return config.FreeEHomeParameterError, ""
	}

	channel, err := strconv.ParseUint(temp[1], 0, 64)
	if err != nil {
		return config.FreeEHomeParameterError, ""
	}

	sequence := atomic.AddUint64(&gw.Sequence, 2)

	reqInviteStreamData := config.ReqInviteStream{
		XMLName:     xml.Name{},
		Version:     2.5,
		Sequence:    sequence,
		CommandType: "REQUEST",
		Method:      "MEDIA",
		Command:     "INVITESTREAM",
		Channel:     channel,
		ChannelType: "SUB",
		SinkIP:      config.StreamStartIP,
		SinkPort:    config.StreamStartPort,
	}

	retBool, err := gw.writeToUDP(reqInviteStreamData, udpClient.UDPAddr, true)
	if retBool {
		//我红方向蓝方发送请求播放实时直播流成功
		logs.BeeLogger.Info("inviteStream request succeeded, deviceID=%s, channelID=%s", deviceID, channelID)
		fmt.Printf("%s inviteStream request succeeded, deviceID=%s,channelID=%s\n", time.Now().Format("2006-01-02 15:04:05"), deviceID, channelID)
		//使用写锁
		gw.seqMutex.Lock()
		gw.SequenceMap[sequence] = &SequenceInfo{
			DeviceID:  deviceID,
			Sequence:  sequence,
			Method:    "MEDIA",
			Command:   "INVITESTREAM",
			ChannelID: channelID,
			ConfigCmd: "",
		}
		gw.seqMutex.Unlock()
		return waitGetSession(deviceID, channelID)
	} else {
		//我红方向蓝方发送请求播放实时直播流失败
		logs.BeeLogger.Error("inviteStream request failed, deviceID=%s", deviceID)
		fmt.Printf("%s inviteStream request failed, deviceID=%s\n", time.Now().Format("2006-01-02 15:04:05"), deviceID)
	}
	return config.FreeEHomeServerError, ""
}

//等待设备返回启动实时直播成功指令，然后取出设备返回的session
func waitGetSession(deviceID, channelID string) (int64, string) {
	idleDelay := time.NewTimer(time.Duration(config.WaitStreamSessionTime) * time.Second)
	defer idleDelay.Stop()

	for {
		select {
		case <-idleDelay.C:
			//超时
			return config.FreeEHomeRequestTimeout, ""
		default:
			//每隔100毫秒请求一次数据
			time.Sleep(100 * time.Millisecond)

			udpClient := filterUDPClientFromDeviceID(deviceID)
			if udpClient == nil {
				return config.FreeEHomeDeviceNotOnline, ""
			}

			udpClient.RLock()
			sessionInfo := udpClient.Sessions[channelID]
			udpClient.RUnlock()

			if sessionInfo != nil && sessionInfo.Session != "" {
				return config.FreeEHomeSuccessOK, sessionInfo.Session
			}
		}
	}
}

//发送关闭实时直播流
func SendByeStream(authCode, deviceID, channelID string) int64 {
	gw.RLock()
	udpClient := gw.UDPClientList[deviceID]
	gw.RUnlock()

	if udpClient == nil {
		return config.FreeEHomeDeviceNotOnline
	}

	udpClient.RLock()
	sessionInfo := udpClient.Sessions[channelID]
	udpClient.RUnlock()

	if sessionInfo == nil {
		return config.FreeEHomeChannelIDNotStreamStart
	}

	sequence := atomic.AddUint64(&gw.Sequence, 2)

	reqByeStreamData := config.ReqByeStream{
		XMLName:     xml.Name{},
		Version:     2.5,
		Sequence:    sequence,
		CommandType: "REQUEST",
		Method:      "MEDIA",
		Command:     "BYESTREAM",
		Session:     sessionInfo.Session,
	}

	retBool, err := gw.writeToUDP(reqByeStreamData, udpClient.UDPAddr, true)
	if retBool {
		//我红方向蓝方发送请求关闭实时直播流成功
		logs.BeeLogger.Info("authCode=%s, byeStream request succeeded, deviceID=%s, channelID=%s", authCode, deviceID, channelID)
		//使用写锁
		gw.seqMutex.Lock()
		gw.SequenceMap[sequence] = &SequenceInfo{
			DeviceID:  deviceID,
			Sequence:  sequence,
			Method:    "MEDIA",
			Command:   "BYESTREAM",
			ChannelID: channelID,
			ConfigCmd: "",
		}
		gw.seqMutex.Unlock()
	} else {
		//我红方向蓝方发送请求关闭实时直播流失败
		logs.BeeLogger.Error("authCode=%s, byeStream request failed, deviceID=%s, channelID=%s, error:%s", authCode, deviceID, channelID, err)
		return config.FreeEHomeServerError
	}

	return config.FreeEHomeSuccessOK
}

//收到Hook关闭信令时删除内存中无效的session
func StreamNotFound(deviceID, channelID string) {
	gw.Lock()
	udpClient := gw.UDPClientList[deviceID]
	if udpClient == nil {
		gw.Unlock()
		return
	}
	delete(udpClient.Sessions, channelID)
	gw.UDPClientList[deviceID] = udpClient
	gw.Unlock()
}
