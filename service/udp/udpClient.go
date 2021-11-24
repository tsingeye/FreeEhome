package udp

import (
	"encoding/xml"
	"fmt"
	"net"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/tsingeye/FreeEhome/config"
	"github.com/tsingeye/FreeEhome/tools"
	"github.com/tsingeye/FreeEhome/tools/logs"
	"github.com/tsingeye/FreeEhome/tools/sqlDB"
)

var gw *GatewayForUDP

type GatewayForUDP struct {
	UDPClientList map[string]*ClientForUDP //键是DeviceID
	Sequence      uint64                   //序列号，用于向设备发送指令，序列号唯一，很关键
	SequenceMap   map[uint64]*SequenceInfo //键是Sequence序列号，用于接收我红方主动发送的信令回复值
	seqMutex      sync.RWMutex             //序列号读写锁
	UDPConn       *net.UDPConn
	sync.RWMutex
}

type ClientForUDP struct {
	DeviceID          string //设备ID
	UDPConn           *net.UDPConn
	UDPAddr           *net.UDPAddr
	KeepaliveCheckNum int8                    //检测设备心跳情况，判断设备是否断开连接
	UDPConnKey        string                  //UDP连接IP和Port，用于检测同一台设备重复注册关闭旧的心跳检测协程
	Sessions          map[string]*SessionInfo //用于记录实时直播生成的session，键是channelID，值是SessionInfo
	sync.RWMutex
}

type SessionInfo struct {
	DeviceID  string
	ChannelID string
	Session   string
}

type SequenceInfo struct {
	DeviceID  string
	Sequence  uint64
	Method    string
	Command   string
	ChannelID string //用于实时直播启动和关闭
	ConfigCmd string //用于区分GetDeviceConfig类型
}

//处理接收到的UDP客户端的数据
func (gateway *GatewayForUDP) readFromUDP() {
	defer gateway.UDPConn.Close()
	for {
		buf := make([]byte, 1024)
		n, udpAddr, err := gateway.UDPConn.ReadFromUDP(buf)
		if err != nil {
			logs.BeeLogger.Error("udpConn.ReadFromUDP() error: %s", err)
			return
		}

		logs.BeeLogger.Debug("UDP原始数据(udp raw data)：%s", string(buf[:n]))

		utf8Data, err := tools.GB2312ToUTF8(buf[:n])
		if err != nil {
			logs.BeeLogger.Error("GB2312ToUTF8() error: %s", err)
			continue
		}

		udpBaseMsg := config.UDPBaseMsg{}
		err = xml.Unmarshal(utf8Data, &udpBaseMsg)
		if err != nil {
			logs.BeeLogger.Error("xml.Unmarshal() error: %s", err)
			continue
		}

		switch udpBaseMsg.CommandType {
		case "REQUEST":
			//接收到的蓝方请求UDP信令
			switch udpBaseMsg.Command {
			case "REGISTER":
				//蓝方向我红方发送注册请求
				go gateway.replyRegisterPacket(utf8Data, udpAddr)
			case "KEEPALIVE":
				//蓝方向我红方发送心跳保活请求
				go gateway.replyKeepalivePacket(utf8Data, udpAddr)
			default:
				logs.BeeLogger.Info("Unknown Command = %s for %s", udpBaseMsg.Command, udpBaseMsg.CommandType)
			}
		case "RESPONSE":
			//接收到的蓝方回复UDP信令
			switch udpBaseMsg.WhichCommand {
			case "GETDEVICECONFIG":
				//接收到蓝方回复的获取设备配置信息，存储到数据库，需注意区分是何种请求
				go gateway.writeDeviceConfigToDB(utf8Data, udpBaseMsg.Sequence, udpAddr)
			case "GETDEVICEWORKSTATUS":
				//接收到蓝方回复的设备工作状态信息，包含当前通道状态，将通道信息写入数据库
				go gateway.writeDeviceWorkStatusToDB(utf8Data, udpBaseMsg.Sequence, udpAddr)
			case "INVITESTREAM":
				//接收到蓝方回复的请求播放实时直播流信令
				go gateway.inviteStreamHandle(utf8Data, udpBaseMsg.Sequence)
			case "BYESTREAM":
				//接收到蓝方回复的请求关闭实时直播信令
				go gateway.byeStreamHandle(utf8Data, udpBaseMsg.Sequence)
			case "QUERYRECORDEDFILES":
				//接收到蓝方回复的请求查询录像文件信令
			case "PTZCONTROL":
				//接收到云台控制的回复信令，忽略不处理
				logs.BeeLogger.Info("ptz response = %s", utf8Data)
			default:
				logs.BeeLogger.Info("Unknown WhichCommand = %s for %s", udpBaseMsg.WhichCommand, udpBaseMsg.CommandType)
			}
		default:
			logs.BeeLogger.Info("Unknown CommandType: %s", udpBaseMsg.CommandType)
		}
	}
}

//发送数据到UDP客户端
func (gateway *GatewayForUDP) writeToUDP(data interface{}, udpAddr *net.UDPAddr, isReq bool) (retBool bool, err error) {
	packetBytes, err := xml.MarshalIndent(data, "  ", "    ")
	if err != nil {
		return
	}

	gb2312Data, err := tools.UTF8ToGB2312(packetBytes)
	if err != nil {
		return
	}

	sendNum := 1
	if isReq {
		//发送请求包，要发送三次，有一次成功即可，默认是回复包，回复一次即可
		sendNum = 3
	}

	//发送三次，有一次成功即可
	var wg sync.WaitGroup
	wg.Add(sendNum)

	for i := 0; i < sendNum; i++ {
		go func() {
			_, err = gateway.UDPConn.WriteToUDP(gb2312Data, udpAddr)
			logs.BeeLogger.Debug("send to UDPClient data is: %s", packetBytes)
			if err == nil {
				retBool = true
			}
			wg.Done()
		}()
	}

	wg.Wait()

	return
}

//回复注册请求
func (gateway *GatewayForUDP) replyRegisterPacket(utf8Data []byte, udpAddr *net.UDPAddr) {
	reqRegisterData := config.ReqRegister{}
	err := xml.Unmarshal(utf8Data, &reqRegisterData)
	if err != nil {
		logs.BeeLogger.Error("replyRegisterPacket() => xml.Unmarshal() error: %s", err)
		return
	}
	logs.BeeLogger.Emergency("Receive register request, DeviceID=%s, DeviceIP=%s", reqRegisterData.DeviceID, udpAddr.IP)
	fmt.Printf("%s Receive register request, DeviceID=%s, DeviceIP=%s\n", time.Now().Format("2006-01-02 15:04:05"), reqRegisterData.DeviceID, udpAddr.IP)
	resRegisterData := config.ResRegister{
		XMLName:            xml.Name{},
		Version:            2.5,
		Sequence:           reqRegisterData.Sequence,
		CommandType:        "RESPONSE",
		Method:             reqRegisterData.Command,
		WhichCommand:       reqRegisterData.Command,
		Status:             200,
		Description:        "OK",
		KeepAliveSeconds:   config.XMLConfigInfo.ParamInfo.KeepAliveSeconds,
		AlarmServerIP:      config.XMLConfigInfo.ParamInfo.AlarmServerIP,
		AlarmServerPort:    config.XMLConfigInfo.ParamInfo.AlarmServerUdpPort,
		AlarmServerType:    config.XMLConfigInfo.ParamInfo.AlarmServerType,
		NTPServerIP:        config.XMLConfigInfo.ParamInfo.NTPServerIP,
		NTPServerPort:      config.XMLConfigInfo.ParamInfo.NTPServerPort,
		NTPInterval:        config.XMLConfigInfo.ParamInfo.NTPInterval,
		PictureServer:      config.XMLConfigInfo.ParamInfo.PictureServerIP,
		PictureServerPort:  config.XMLConfigInfo.ParamInfo.PictureServerPort,
		PictureServerType:  config.XMLConfigInfo.ParamInfo.PictureServerType,
		BlackListAddr:      "",
		BlackListName:      "",
		BlackListPort:      0,
		BlackListUser:      "",
		BlackListPassword:  "",
		RegisterServerIP:   "",
		RegisterServerPort: 0,
	}

	retBool, err := gateway.writeToUDP(resRegisterData, udpAddr, false)
	if retBool {
		//蓝方向我红方进行设备注册请求回复成功
		logs.BeeLogger.Info("register request reply succeeded, deviceID=%s", reqRegisterData.DeviceID)
	} else {
		//注册请求回复失败
		logs.BeeLogger.Error("register request reply failed, deviceID=%s, error: %s", reqRegisterData.DeviceID, err)
	}

	udpConnKey := tools.StringsJoin(udpAddr.IP.String(), ":", strconv.Itoa(udpAddr.Port))

	//注意此处要使用写锁，目的是只在连续多个注册请求中获取GetServerInfo配置信息和心跳协程只开启一次
	gateway.Lock()
	udpClient, ok := gateway.UDPClientList[reqRegisterData.DeviceID]
	defer gateway.Unlock()

	if !ok || udpClient.UDPConnKey != udpConnKey {
		client := &ClientForUDP{
			DeviceID:          reqRegisterData.DeviceID,
			UDPConn:           gateway.UDPConn,
			UDPAddr:           udpAddr,
			KeepaliveCheckNum: 0,
			Sessions:          make(map[string]*SessionInfo),
			UDPConnKey:        udpConnKey,
		}

		//存储已注册成功的设备信息，使用写锁
		gateway.UDPClientList[reqRegisterData.DeviceID] = client
		//存储设备列表
		go func() {
			deviceList := &sqlDB.DeviceList{
				DeviceID:     reqRegisterData.DeviceID,
				DeviceIP:     udpAddr.IP.String(),
				DeviceName:   "",
				SerialNumber: "",
				Status:       "ON",
			}

			if sqlDB.Save(deviceList) {
				tabName := sqlDB.GetTableName(&sqlDB.DeviceList{})
				logs.BeeLogger.Info("deviceID=%s saved record to %s's table successfully", reqRegisterData.DeviceID, tabName)
			}
		}()
		//开协程检测心跳保活
		go gateway.keepaliveCheck(reqRegisterData.DeviceID, udpConnKey)
		//获取GetServerInfo配置信息
		go gateway.sendGetServerInfoPacket(reqRegisterData.DeviceID, udpAddr)
		//获取GetDeviceWorkStatus状态信息
		go gateway.sendGetDeviceWorkStatusPacket(reqRegisterData.DeviceID, udpAddr)
	}
}

//心跳检测，判断设备是否断线
func (gateway *GatewayForUDP) keepaliveCheck(deviceID, udpConnKey string) {
	t1 := time.NewTicker(time.Duration(config.XMLConfigInfo.ParamInfo.KeepAliveSeconds) * time.Second)
	defer t1.Stop()
	for {
		select {
		case <-t1.C:
			gateway.RLock()
			udpClient := gateway.UDPClientList[deviceID]
			gateway.RUnlock()

			if udpClient == nil {
				//设备断开连接
				go func() {
					//更新设备列表，将Online置为on
					sqlDB.UpdateTable(fmt.Sprintf(`UPDATE %s SET Status = 'OFF', UpdatedAt= '%v' WHERE DeviceID= '%s'`, sqlDB.GetTableName(&sqlDB.DeviceList{}), time.Now(), deviceID))
				}()

				go func() {
					//更新设备通道列表，将Status置为OFF
					sqlDB.UpdateTable(fmt.Sprintf(`UPDATE %s SET Status = 'OFF', UpdatedAt= '%v' WHERE DeviceID= '%s'`, sqlDB.GetTableName(&sqlDB.ChannelList{}), time.Now(), deviceID))
				}()

				logs.BeeLogger.Error("deviceID=%s disconnect", deviceID)
				fmt.Printf("%s deviceID=%s disconnect\n", time.Now().Format("2006-01-02 15:04:05"), deviceID)
				return
			} else if udpClient.UDPConnKey != udpConnKey {
				//同一个设备再次注册成功，此时需关闭旧的心跳检测协程
				fmt.Printf("%s DeviceID=%s reconnected successfully, closed the old keepalive goroutine\n", time.Now().Format("2006-01-02 15:04:05"), deviceID)
				logs.BeeLogger.Warn("DeviceID=%s reconnected successfully, closed the old keepalive goroutine", deviceID)
				return
			} else {
				if udpClient.KeepaliveCheckNum+1 >= 6 {
					fmt.Printf("%s DeviceID=%s timeout disconnect, please re-register\n", time.Now().Format("2006-01-02 15:04:05"), deviceID)
					logs.BeeLogger.Warn("DeviceID=%s timeout disconnect, please re-register", deviceID)
					//超时，设备断开连接，内存中删除设备信息，使用写锁
					gateway.Lock()
					delete(gateway.UDPClientList, deviceID)
					gateway.Unlock()
					//设备断开连接
					go func() {
						//更新设备列表，将Status置为OFF
						sqlDB.UpdateTable(fmt.Sprintf(`UPDATE %s SET Status = 'OFF', UpdatedAt= '%v' WHERE DeviceID= '%s'`, sqlDB.GetTableName(&sqlDB.DeviceList{}), time.Now(), deviceID))
					}()

					go func() {
						//更新设备通道列表，将Status置为OFF
						sqlDB.UpdateTable(fmt.Sprintf(`UPDATE %s SET Status = 'OFF', UpdatedAt= '%v' WHERE DeviceID= '%s'`, sqlDB.GetTableName(&sqlDB.ChannelList{}), time.Now(), deviceID))
					}()

					return
				} else {
					//使用写锁
					gateway.Lock()
					udpClient.KeepaliveCheckNum++
					gateway.UDPClientList[udpClient.DeviceID] = udpClient
					gateway.Unlock()
				}
			}

		}
	}
}

//回复心跳保活
func (gateway *GatewayForUDP) replyKeepalivePacket(utf8Data []byte, udpAddr *net.UDPAddr) {
	reqKeepaliveData := config.ReqKeepalive{}
	err := xml.Unmarshal(utf8Data, &reqKeepaliveData)
	if err != nil {
		logs.BeeLogger.Error("replyKeepalivePacket() => xml.Unmarshal() error: %s", err)
		return
	}
	//读锁
	gateway.RLock()
	udpClient := gateway.UDPClientList[reqKeepaliveData.DeviceID]
	gateway.RUnlock()

	resKeepaliveData := config.ResKeepalive{
		XMLName:      xml.Name{},
		Version:      2.5,
		Sequence:     reqKeepaliveData.Sequence,
		CommandType:  "RESPONSE",
		Method:       "CONTROL",
		WhichCommand: reqKeepaliveData.Command,
		Status:       200,
		Description:  "OK",
	}

	if udpClient == nil {
		//设备未连接，此时若收到心跳包需回复403让设备尽快重新注册
		resKeepaliveData.Method = "REGISTER"
		resKeepaliveData.Status = 403
		logs.BeeLogger.Error("device disconnected, unable to keepalive, please register!")
	} else {
		gateway.Lock()
		//重置心跳保活计数器
		udpClient.KeepaliveCheckNum = 0
		gateway.UDPClientList[udpClient.DeviceID] = udpClient
		gateway.Unlock()
	}

	retBool, err := gateway.writeToUDP(resKeepaliveData, udpAddr, false)
	if !retBool {
		logs.BeeLogger.Error("keepalive request reply failed, deviceID=%s", udpClient.DeviceID)
	}
}

//发送请求，获取GetServerInfo配置信息
func (gateway *GatewayForUDP) sendGetServerInfoPacket(deviceID string, udpAddr *net.UDPAddr) {
	//原子操作，这样使用为了下面存储数据
	sequence := atomic.AddUint64(&gateway.Sequence, 2)

	reqGetServerInfoData := config.ReqGetServerInfo{
		XMLName:      xml.Name{},
		Version:      2.5,
		Sequence:     sequence,
		CommandType:  "REQUEST",
		Method:       "PARAM",
		Command:      "GETDEVICECONFIG",
		ConfigCmd:    "GetServerInfo",
		ConfigParam1: 0,
		ConfigParam2: 0,
		ConfigParam3: 0,
		ConfigParam4: 0,
	}

	retBool, err := gateway.writeToUDP(reqGetServerInfoData, udpAddr, true)
	if retBool {
		//我红方向蓝方发送获取GetServerInfo信令成功
		logs.BeeLogger.Info("getServerInfo request succeeded, deviceID=%s", deviceID)
		//使用写锁
		gateway.seqMutex.Lock()
		defer gateway.seqMutex.Unlock()
		gateway.SequenceMap[sequence] = &SequenceInfo{
			DeviceID:  deviceID,
			Sequence:  sequence,
			Method:    "PARAM",
			Command:   "GETDEVICECONFIG",
			ChannelID: "",
			ConfigCmd: "GetServerInfo",
		}
	} else {
		//我红方向蓝方发送获取GetServerInfo信令失败
		logs.BeeLogger.Error("getServerInfo request failed, error: %s", err)
	}
}

//发送请求，获取GetDeviceWorkStatus状态信息
func (gateway *GatewayForUDP) sendGetDeviceWorkStatusPacket(deviceID string, udpAddr *net.UDPAddr) {
	time.Sleep(60 * time.Second)
	//原子操作，这样使用为了下面存储数据
	sequence := atomic.AddUint64(&gateway.Sequence, 2)

	reqGetDeviceWorkStatus := config.ReqGetDeviceWorkStatus{
		XMLName:     xml.Name{},
		Version:     4.0,
		Sequence:    sequence,
		CommandType: "REQUEST",
		Method:      "QUERY",
		Command:     "GETDEVICEWORKSTATUS",
		Params:      struct{}{},
	}

	retBool, err := gateway.writeToUDP(reqGetDeviceWorkStatus, udpAddr, false)
	if retBool {
		//我红方向蓝方发送获取GetDeviceWorkStatus信令成功
		logs.BeeLogger.Info("getDeviceWorkStatus request succeeded, deviceID=%s", deviceID)
		//使用写锁
		gateway.seqMutex.Lock()
		defer gateway.seqMutex.Unlock()
		gateway.SequenceMap[sequence] = &SequenceInfo{
			DeviceID:  deviceID,
			Sequence:  sequence,
			Method:    "QUERY",
			Command:   "GETDEVICEWORKSTATUS",
			ChannelID: "",
			ConfigCmd: "",
		}
	} else {
		//我红方向蓝方发送获取GetDeviceWorkStatus信令失败
		logs.BeeLogger.Error("getDeviceWorkStatus request failed, error: %s", err)
	}
}

//将获得的设备配置信息写入数据库
func (gateway *GatewayForUDP) writeDeviceConfigToDB(utf8Data []byte, sequence uint64, udpAddr *net.UDPAddr) {
	gateway.seqMutex.Lock()
	sequenceInfo, ok := gateway.SequenceMap[sequence]
	if !ok {
		//解锁
		gateway.seqMutex.Unlock()
		return
	}
	logs.BeeLogger.Debug("write deviceConfig to DB")
	//删除数据，保证每次请求的设备配置信息只写入数据库一次
	delete(gateway.SequenceMap, sequence)
	gateway.seqMutex.Unlock()

	//处理GetDeviceConfig中具体请求回复的内容
	switch sequenceInfo.ConfigCmd {
	case "GetServerInfo":
		writeServerInfoToDB(sequenceInfo.DeviceID, utf8Data, udpAddr)
	}
}

//处理获得是设备工作状态信息
func (gateway *GatewayForUDP) writeDeviceWorkStatusToDB(utf8Data []byte, sequence uint64, udpAddr *net.UDPAddr) {
	gateway.seqMutex.Lock()
	sequenceInfo, ok := gateway.SequenceMap[sequence]
	if !ok {
		//解锁
		gateway.seqMutex.Unlock()
		return
	}
	logs.BeeLogger.Debug("write deviceConfig to DB")
	//删除数据，保证每次请求的设备状态信息只写入数据库一次
	delete(gateway.SequenceMap, sequence)
	gateway.seqMutex.Unlock()
	fmt.Printf("%s get DeviceWorkStatus, DeviceID=%s\n", time.Now().Format("2006-01-02 15:04:05"), sequenceInfo.DeviceID)
	writeGetDeviceWorkStatusToDB(sequenceInfo.DeviceID, utf8Data, udpAddr)
}

//处理获得的实时直播返回内容
func (gateway *GatewayForUDP) inviteStreamHandle(utf8Data []byte, sequence uint64) {
	resInviteStreamData := config.ResInviteStream{}
	err := xml.Unmarshal(utf8Data, &resInviteStreamData)
	if err != nil {
		logs.BeeLogger.Error("inviteStreamHandle() => xml.Unmarshal() error: %s", err)
		return
	}

	gateway.seqMutex.RLock()
	sequenceInfo := gateway.SequenceMap[sequence]
	gateway.seqMutex.RUnlock()
	if sequenceInfo == nil {
		logs.BeeLogger.Error("inviteStreamHandle() ==> SequenceMap[%s] is nil, please check it", sequence)
		return
	}
	if resInviteStreamData.Status == 200 {
		logs.BeeLogger.Info("deviceID=%s, channelID=%s inviteStream start success, get session=%s", sequenceInfo.DeviceID, sequenceInfo.ChannelID, resInviteStreamData.Session)
		fmt.Printf("%s deviceID=%s, channelID=%s inviteStream start success, get session=%s\n", time.Now().Format("2006-01-02 15:04:05"), sequenceInfo.DeviceID, sequenceInfo.ChannelID, resInviteStreamData.Session)

		gateway.Lock()
		udpClient := gateway.UDPClientList[sequenceInfo.DeviceID]
		if udpClient == nil {
			logs.BeeLogger.Error("get sessionURL success, deviceID=%s, channelID=%s, but UDPClientList[%s] is nil, please check it", sequenceInfo.DeviceID, sequenceInfo.ChannelID, sequenceInfo.DeviceID)
			gateway.Unlock()
			return
		}

		udpClient.Sessions[sequenceInfo.ChannelID] = &SessionInfo{
			DeviceID:  sequenceInfo.DeviceID,
			ChannelID: sequenceInfo.ChannelID,
			Session:   resInviteStreamData.Session,
		}
		gateway.UDPClientList[sequenceInfo.DeviceID] = udpClient
		gateway.Unlock()
	} else {
		logs.BeeLogger.Info("deviceID=%s inviteStream start failed, status=%d", sequenceInfo.DeviceID, resInviteStreamData.Status)
		fmt.Printf("%s deviceID=%s inviteStream start failed, status=%d\n", time.Now().Format("2006-01-02 15:04:05"), sequenceInfo.DeviceID, resInviteStreamData.Status)
	}
}

//处理关闭实时直播返回内容
func (gateway *GatewayForUDP) byeStreamHandle(utf8Data []byte, sequence uint64) {
	resByeStreamData := config.ResByeStream{}
	err := xml.Unmarshal(utf8Data, &resByeStreamData)
	if err != nil {
		logs.BeeLogger.Error("byeStreamHandle() => xml.Unmarshal() error: %s", err)
		return
	}
	if resByeStreamData.Status == 200 {
		logs.BeeLogger.Warn("session=%s has been closed!", resByeStreamData.Session)
		fmt.Printf("%s session=%s has been closed!\n", time.Now().Format("2006-01-02 15:04:05"), resByeStreamData.Session)

		gateway.seqMutex.RLock()
		sequenceInfo := gateway.SequenceMap[sequence]
		gateway.seqMutex.RUnlock()
		if sequenceInfo == nil {
			logs.BeeLogger.Error("byeStreamHandle() ==> SequenceMap[%s] is nil, please check it", sequence)
			return
		}

		gateway.Lock()
		udpClient := gateway.UDPClientList[sequenceInfo.DeviceID]
		if udpClient == nil {
			gateway.Unlock()
			return
		}
		delete(udpClient.Sessions, sequenceInfo.ChannelID)
		gateway.UDPClientList[sequenceInfo.DeviceID] = udpClient
		gateway.Unlock()
	}
}
