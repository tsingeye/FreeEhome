package udp

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/tsingeye/FreeEhome/config"
	"github.com/tsingeye/FreeEhome/tools"
	"github.com/tsingeye/FreeEhome/tools/logs"
	"io/ioutil"
	"net"
	"net/http"
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

	udpClient.RLock()
	sessionInfo, ok := udpClient.Sessions[channelID]
	udpClient.RUnlock()

	if ok && sessionInfo.Session != "" {
		fmt.Printf("%s authCode=%s, DeviceID=%s, ChannelID=%s, this inviteStream's session in memory, direct request for STS Service\n", time.Now().Format("2006-01-02 15:04:05"), authCode, deviceID, channelID)
		logs.BeeLogger.Emergency("authCode=%s, DeviceID=%s, ChannelID=%s, this inviteStream's session in memory, there is no need to get a new session", authCode, deviceID, channelID)
		//该设备对应的通道已被触发过实时直播，此时无需向设备发送启动指令，直接向STS服务器发送POST请求
		sessionURL := getSessionURLFromSTS(sessionInfo.Session)
		if sessionURL == "" {
			//当sessionURL为空即可认为STS服务出现问题，删除已保存的数据
			udpClient.Lock()
			delete(udpClient.Sessions, channelID)
			udpClient.Unlock()

			gw.Lock()
			gw.UDPClientList[deviceID] = udpClient
			gw.Unlock()

			return config.FreeEHomeRequestTimeout, ""
		} else {
			return config.FreeEHomeSuccessOK, sessionURL
		}
	}

	fmt.Printf("%s DeviceID=%s, ChannelID=%s, this inviteStream's session not in memory, please send inviteStream, after request STS Service\n", time.Now().Format("2006-01-02 15:04:05"), deviceID, channelID)
	logs.BeeLogger.Emergency("DeviceID=%s, ChannelID=%s, this inviteStream's session not in memory, please send inviteStream, after request STS Service", deviceID, channelID)
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
		logs.BeeLogger.Info("inviteStream request succeeded, deviceID=%s", deviceID)
		fmt.Printf("%s inviteStream request succeeded, deviceID=%s\n", time.Now().Format("2006-01-02 15:04:05"), deviceID)
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
		return waitGetSessionURL(deviceID, channelID)
	} else {
		//我红方向蓝方发送请求播放实时直播流失败
		logs.BeeLogger.Error("inviteStream request failed, deviceID=%s", deviceID)
		fmt.Printf("%s inviteStream request failed, deviceID=%s\n", time.Now().Format("2006-01-02 15:04:05"), deviceID)
	}
	return config.FreeEHomeServerError, ""
}

//获取实时视频播放URL，第一次内存中无此URL，向设备发送start指令，等设备返回status=200时向STS服务发送POST请求URL
//当内存中有URL时，直接向STS服务发送POST请求URL，不直接返回是防止URL地址因为STS服务器断开等原因而失效
func getSessionURLFromSTS(session string) string {
	client := &http.Client{
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				conn, err := net.DialTimeout(netw, addr, time.Second*2) //设置建立连接超时
				if err != nil {
					return nil, err
				}
				conn.SetDeadline(time.Now().Add(time.Second * 3)) //设置发送接受数据超时
				return conn, nil
			},
			ResponseHeaderTimeout: time.Second * 4,
		},
	}
	data := struct {
		SessionID string `json:"sessionId"`
	}{
		SessionID: session,
	}

	sessionByte, err := json.MarshalIndent(data, "  ", "    ")
	if err != nil {
		logs.BeeLogger.Error("getSessionURLFromSTS() => xml.MarshalIndent() error: %s", err)
		return ""
	}

	reqest, err := http.NewRequest("POST", tools.StringsJoin("http://", config.STSAddr, "/eh/StartStream"), bytes.NewReader(sessionByte))
	if err != nil {
		logs.BeeLogger.Error("getSessionURLFromSTS() => http.NewRequest() error: %s", err)
		return ""
	}

	//增加Header选项
	reqest.Header.Set("Content-Type", "application/json")

	response, err := client.Do(reqest)
	if err != nil {
		logs.BeeLogger.Error("getSessionURLFromSTS() => http.Do() error: %s", err)
		return ""
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		logs.BeeLogger.Error("getSessionURLFromSTS() => ioutil.ReadAll() error: %s", err)
		return ""
	}
	logs.BeeLogger.Debug("getSessionURLFromSTS() body is: %s", string(body))
	bodyData := struct {
		ErrorCode    int64  `json:"errorCode"`
		ErrorCodeStr string `json:"errorCodeStr"`
		PushUrl      string `json:"pushUrl"`
	}{}
	err = json.Unmarshal(body, &bodyData)
	if err != nil {
		logs.BeeLogger.Error("getSessionURLFromSTS() => json.Unmarshal() error: %s", err)
		return ""
	}
	if bodyData.ErrorCode == 0 {
		return bodyData.PushUrl
	}

	return ""
}

//等待设备返回启动实时直播成功指令，然后从内存取STS返回的sessionURL
func waitGetSessionURL(deviceID, channelID string) (int64, string) {
	idleDelay := time.NewTimer(time.Duration(config.WaitStreamURLTime) * time.Second)
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

			if sessionInfo != nil && sessionInfo.SessionURL != "" {
				return config.FreeEHomeSuccessOK, sessionInfo.SessionURL
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
