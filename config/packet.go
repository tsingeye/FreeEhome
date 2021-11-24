package config

import "encoding/xml"

//********************************TCP客户端数据结构********************************
//TCP客户端通信数据结构体
type TCPClientData struct {
	IVMS IVMSData `json:"iVMS"`
}

type IVMSData struct {
	Body   interface{} `json:"body"`
	Header IVMSHeader  `json:"header"`
}

type IVMSHeader struct {
	Version  string `json:"version"`
	Protocol string `json:"protocol"`
	Csq      string `json:"csq"`
}

//********************************websocket数据结构********************************
//websocket协议数据格式
type WSBaseMsg struct {
	Body   interface{} `json:"body"`
	Header WSHeader    `json:"header"`
}

type WSHeader struct {
	Protocol string `json:"protocol"` //协议
}

//用于解析websocket客户端请求设备列表
type WSDeviceList struct {
	Body   WSDeviceListBody `json:"body"`
	Header WSHeader         `json:"header"`
}

type WSDeviceListBody struct {
	Start uint64 `json:"start"` //分页开始值，默认从0开始
	Limit uint64 `json:"limit"` //分页大小，默认查询10个
}

//用于解析websocket客户端请求设备通道列表
type WSChannelList struct {
	Body   WSChannelListBody `json:"body"`
	Header WSHeader          `json:"header"`
}

type WSChannelListBody struct {
	Start uint64 `json:"start"` //分页开始值，默认从0开始
	Limit uint64 `json:"limit"` //分页大小，默认查询10个
}

//解析实时直播开启和关闭信令数据
type WSStreamStartOrStop struct {
	Body   WSStreamStartOrStopBody `json:"body"`
	Header WSHeader                `json:"header"`
}

type WSStreamStartOrStopBody struct {
	ChannelID  string `json:"channelID"`
	StreamType string `json:"streamType"`
}

//********************************配置文件解析config.xml********************************
//用于解析配置文件中的XML文件
type XMLConfig struct {
	XMLName        xml.Name `xml:"LocalCfg"`
	ParamInfo      param
	AddressMapInfo addressMap
	DasInfoData    dasInfo
}

type param struct {
	XMLName               xml.Name `xml:"Param"`
	LogPath               string   `xml:"LogPath"`
	LogLevel              uint64   `xml:"LogLevel"`
	LogAutoDel            bool     `xml:"LogAutoDel"`
	KeepAliveSeconds      uint64   `xml:"KeepAliveSeconds"`
	KeepAliveCount        uint64   `xml:"KeepAliveCount"`
	AlarmServerType       uint64   `xml:"AlarmServerType"`
	AlarmServerIP         string   `xml:"AlarmServerIP"`
	AlarmServerUdpPort    uint64   `xml:"AlarmServerUdpPort"`
	AlarmServerTcpPort    uint64   `xml:"AlarmServerTcpPort"`
	AlarmServerMqttPort   uint64   `xml:"AlarmServerMqttPort"`
	AlarmServerPortUseCms uint64   `xml:"AlarmServerPortUseCms"`
	NTPServerIP           string   `xml:"NTPServerIP"`
	NTPServerPort         uint64   `xml:"NTPServerPort"`
	NTPInterval           uint64   `xml:"NTPInterval"`
	PictureServerType     uint64   `xml:"PictureServerType"`
	PictureServerIP       string   `xml:"PictureServerIP"`
	PictureServerPort     uint64   `xml:"PictureServerPort"`
	CmsAccessSecurity     uint64   `xml:"CmsAccessSecurity"`
	AlarmAccessSecurity   uint64   `xml:"AlarmAccessSecurity"`
	StreamAccessSecurity  uint64   `xml:"StreamAccessSecurity"`
	AlarmPictureSeparate  uint64   `xml:"AlarmPictureSeparate"`
}

type addressMap struct {
	XMLName             xml.Name `xml:"AddressMap"`
	Enable              uint64   `xml:"Enable"`
	StreamServerIP      string   `xml:"StreamServerIP"`
	StreamServerPort    uint64   `xml:"StreamServerPort"`
	AudioServerIP       string   `xml:"AudioServerIP"`
	AudioServerPort     uint64   `xml:"AudioServerPort"`
	AlarmServerIP       string   `xml:"AlarmServerIP"`
	AlarmServerUdpPort  uint64   `xml:"AlarmServerUdpPort"`
	AlarmServerTcpPort  uint64   `xml:"AlarmServerTcpPort"`
	AlarmServerMqttPort uint64   `xml:"AlarmServerMqttPort"`
	PictureServerIP     string   `xml:"PictureServerIP"`
	PictureServerPort   uint64   `xml:"PictureServerPort"`
}

type dasInfo struct {
	XMLName     xml.Name `xml:"DasInfo"`
	ServerType  uint64   `xml:"ServerType"`
	DasInfoPath string   `xml:"DasInfoPath"`
	DASIP       string   `xml:"DASIP"`
	DASPort     uint64   `xml:"DASPort"`
}

//********************************UDP通信数据结构********************************
//UDP通信数据结构体
type UDPBaseMsg struct {
	XMLName      xml.Name `xml:"PPVSPMessage"`
	Version      float64  `xml:"Version"`
	Sequence     uint64   `xml:"Sequence"`
	CommandType  string   `xml:"CommandType"`  //REQUEST Or RESPONSE
	Command      string   `xml:"Command"`      //用于解析蓝方请求信令
	WhichCommand string   `xml:"WhichCommand"` //用于解析蓝方回复信令
}

//用于解析蓝方向我红方发送的注册请求UDP信令
type ReqRegister struct {
	XMLName         xml.Name `xml:"PPVSPMessage"`
	Version         float64  `xml:"Version"`
	Sequence        uint64   `xml:"Sequence"`
	CommandType     string   `xml:"CommandType"`
	Command         string   `xml:"Command"`
	NetUnitType     string   `xml:"Params>NetUnitType"`
	DeviceID        string   `xml:"Params>DeviceID"`
	FirmwareVersion string   `xml:"Params>FirmwareVersion"`
	LocalIP         string   `xml:"Params>LocalIP"`
	LocalPort       uint64   `xml:"Params>LocalPort"`
	DevType         uint64   `xml:"Params>DevType"`
	Manufacture     uint64   `xml:"Params>Manufacture"`
}

//用于构造我红方向蓝方回复注册请求回复UDP信令
type ResRegister struct {
	XMLName            xml.Name `xml:"PPVSPMessage"`
	Version            float64  `xml:"Version"`
	Sequence           uint64   `xml:"Sequence"`
	CommandType        string   `xml:"CommandType"`
	Method             string   `xml:"Method"`
	WhichCommand       string   `xml:"WhichCommand"`
	Status             uint64   `xml:"Status"`
	Description        string   `xml:"Description"`
	KeepAliveSeconds   uint64   `xml:"Params>KeepAliveSeconds"`
	AlarmServerIP      string   `xml:"Params>AlarmServerIP"`
	AlarmServerPort    uint64   `xml:"Params>AlarmServerPort"`
	AlarmServerType    uint64   `xml:"Params>AlarmServerType"`
	NTPServerIP        string   `xml:"Params>NTPServerIP"`
	NTPServerPort      uint64   `xml:"Params>NTPServerPort"`
	NTPInterval        uint64   `xml:"Params>NTPInterval"`
	PictureServer      string   `xml:"Params>PictureServer"`
	PictureServerPort  uint64   `xml:"Params>PictureServerPort"`
	PictureServerType  uint64   `xml:"Params>PictureServerType"`
	BlackListAddr      string   `xml:"Params>BlackListAddr"`
	BlackListName      string   `xml:"Params>BlackListName"`
	BlackListPort      uint64   `xml:"Params>BlackListPort"`
	BlackListUser      string   `xml:"Params>BlackListUser"`
	BlackListPassword  string   `xml:"Params>BlackListPassword"`
	RegisterServerIP   string   `xml:"Params>RegisterServerIP"`
	RegisterServerPort uint64   `xml:"Params>RegisterServerPort"`
}

//用于解析蓝方向我红方发送的心跳保活请求UDP信令
type ReqKeepalive struct {
	XMLName     xml.Name `xml:"PPVSPMessage"`
	Version     float64  `xml:"Version"`
	Sequence    uint64   `xml:"Sequence"`
	CommandType string   `xml:"CommandType"`
	Command     string   `xml:"Command"`
	DeviceID    string   `xml:"Params>DeviceID"`
}

//用于构造我红方向蓝方回复心跳保活请求回复UDP信令
type ResKeepalive struct {
	XMLName      xml.Name `xml:"PPVSPMessage"`
	Version      float64  `xml:"Version"`
	Sequence     uint64   `xml:"Sequence"`
	CommandType  string   `xml:"CommandType"`
	Method       string   `xml:"Method"`
	WhichCommand string   `xml:"WhichCommand"`
	Status       uint64   `xml:"Status"`
	Description  string   `xml:"Description"`
}

//用于构造我红方向蓝方发送获取GetServerInfo请求信令
type ReqGetServerInfo struct {
	XMLName      xml.Name `xml:"PPVSPMessage"`
	Version      float64  `xml:"Version"`
	Sequence     uint64   `xml:"Sequence"`
	CommandType  string   `xml:"CommandType"`
	Method       string   `xml:"Method"`
	Command      string   `xml:"Command"`
	ConfigCmd    string   `xml:"Params>ConfigCmd"`
	ConfigParam1 uint64   `xml:"Params>ConfigParam1"`
	ConfigParam2 uint64   `xml:"Params>ConfigParam2"`
	ConfigParam3 uint64   `xml:"Params>ConfigParam3"`
	ConfigParam4 uint64   `xml:"Params>ConfigParam4"`
}

//用于解析蓝方向我红方回复的GetServerInfo回复UDP信令
type ResGetServerInfo struct {
	XMLName            xml.Name `xml:"PPVSPMessage"`
	Version            float64  `xml:"Version"`
	Sequence           uint64   `xml:"Sequence"`
	CommandType        string   `xml:"CommandType"`
	WhichCommand       string   `xml:"WhichCommand"`
	Status             uint64   `xml:"Status"`
	Description        string   `xml:"Description"`
	ChannelNumber      uint64   `xml:"Params>ConfigXML>SERVERINFO>ChannelNumber"`
	ChannelAmount      uint64   `xml:"Params>ConfigXML>SERVERINFO>ChannelAmount"`
	DVRType            uint64   `xml:"Params>ConfigXML>SERVERINFO>DVRType"`
	DiskNumber         uint64   `xml:"Params>ConfigXML>SERVERINFO>DiskNumber"`
	SerialNumber       string   `xml:"Params>ConfigXML>SERVERINFO>SerialNumber"`
	AlarmInPortNum     uint64   `xml:"Params>ConfigXML>SERVERINFO>AlarmInPortNum"`
	AlarmInAmount      uint64   `xml:"Params>ConfigXML>SERVERINFO>AlarmInAmount"`
	AlarmOutPortNum    uint64   `xml:"Params>ConfigXML>SERVERINFO>AlarmOutPortNum"`
	AlarmOutAmount     uint64   `xml:"Params>ConfigXML>SERVERINFO>AlarmOutAmount"`
	StartChannel       uint64   `xml:"Params>ConfigXML>SERVERINFO>StartChannel"`
	AudioChanNum       uint64   `xml:"Params>ConfigXML>SERVERINFO>AudioChanNum"`
	MaxDigitChannelNum uint64   `xml:"Params>ConfigXML>SERVERINFO>MaxDigitChannelNum"`
	AudioEncType       uint64   `xml:"Params>ConfigXML>SERVERINFO>AudioEncType"`
	SmartType          uint64   `xml:"Params>ConfigXML>SERVERINFO>SmartType"`
	StartChan          uint64   `xml:"Params>ConfigXML>SERVERINFO>StartChan"`
	StartDChan         uint64   `xml:"Params>ConfigXML>SERVERINFO>StartDChan"`
}

//用于构造我红方向蓝方发送获取设备状态GetDeviceWorkStatus请求信令
type ReqGetDeviceWorkStatus struct {
	XMLName     xml.Name `xml:"PPVSPMessage"`
	Version     float64  `xml:"Version"`
	Sequence    uint64   `xml:"Sequence"`
	CommandType string   `xml:"CommandType"`
	Method      string   `xml:"Method"`
	Command     string   `xml:"Command"`
	Params      struct{} `xml:"Params"`
}

//用于解析蓝方向我红方回复的GetDeviceWorkStatus UDP信令
type ResGetDeviceWorkStatus struct {
	XMLName            xml.Name `xml:"PPVSPMessage"`
	Version            float64  `xml:"Version"`
	Sequence           uint64   `xml:"Sequence"`
	CommandType        string   `xml:"CommandType"`
	WhichCommand       string   `xml:"WhichCommand"`
	Status             string   `xml:"Status"`
	Description        string   `xml:"Description"`
	Run                uint64   `xml:"Params>DeviceStatusXML>Run"`                //设备状态：0-正常，1-CPU使用率过高（高于85%），2-硬件错误
	CPU                uint64   `xml:"Params>DeviceStatusXML>CPU"`                //CPU使用率，取值范围为0%到100%
	Mem                uint64   `xml:"Params>DeviceStatusXML>Mem"`                //内存使用率，取值范围为0%到100%
	DSKStatus          dskList  `xml:"Params>DeviceStatusXML>DSKStatus"`          //包含多个参数的字符串，包括硬盘号，硬盘容量（单位：MB），硬盘剩余空间（单位：MB），和硬盘状态（0-活跃，1-睡眠，2-异常）；每个参数由“-”分隔
	CHStatus           chList   `xml:"Params>DeviceStatusXML>CHStatus"`           //包含多个参数的字符串，包括通道号，录像状态（0-停止，1-开始），视频信号状态（0-正常，1-视频丢失），通道编码状态（0-正常，1-异常），实际码率（单位：Kbps），关联客户端数量
	AlarmInStatus      string   `xml:"Params>DeviceStatusXML>AlarmInStatus"`      //已启用报警输入号，多个号之间用逗号隔开
	AlarmOutStatus     string   `xml:"Params>DeviceStatusXML>AlarmOutStatus"`     //已启用报警输出号，多个号之间用逗号隔开
	LocalDisplayStatus uint64   `xml:"Params>DeviceStatusXML>LocalDisplayStatus"` //本地显示状态：0-正常，1-异常
	Remark             string   `xml:"Params>DeviceStatusXML>Remark"`
}
type chList struct {
	CHList []string `xml:"CH"`
}

type dskList struct {
	DSKList []string `xml:"DSK"`
}

//用于构造我红方向蓝方发送请求播放实时直播流信令
type ReqInviteStream struct {
	XMLName     xml.Name `xml:"PPVSPMessage"`
	Version     float64  `xml:"Version"`
	Sequence    uint64   `xml:"Sequence"`
	CommandType string   `xml:"CommandType"`
	Method      string   `xml:"Method"`
	Command     string   `xml:"Command"`
	Channel     uint64   `xml:"Params>Channel"`
	ChannelType string   `xml:"Params>ChannelType"`
	SinkIP      string   `xml:"Params>SinkIP"`
	SinkPort    int64    `xml:"Params>SinkPort"`
}

//用于解析蓝方向我红方回复的请求播放实时直播流UDP信令
type ResInviteStream struct {
	XMLName      xml.Name `xml:"PPVSPMessage"`
	Version      float64  `xml:"Version"`
	Sequence     int64    `xml:"Sequence"`
	CommandType  string   `xml:"CommandType"`
	WhichCommand string   `xml:"WhichCommand"`
	Status       int64    `xml:"Status"`
	Description  string   `xml:"Description"`
	Session      string   `xml:"Params>Session"`
}

//用于构造我红方向蓝方发送请求关闭实时直播流信令
type ReqByeStream struct {
	XMLName     xml.Name `xml:"PPVSPMessage"`
	Version     float64  `xml:"Version"`
	Sequence    uint64   `xml:"Sequence"`
	CommandType string   `xml:"CommandType"`
	Method      string   `xml:"Method"`
	Command     string   `xml:"Command"`
	Session     string   `xml:"Params>Session"`
}

//用于解析蓝方向我红方回复的请求关闭实时直播流UDP信令
type ResByeStream struct {
	XMLName      xml.Name `xml:"PPVSPMessage"`
	Version      float64  `xml:"Version"`
	Sequence     int64    `xml:"Sequence"`
	CommandType  string   `xml:"CommandType"`
	WhichCommand string   `xml:"WhichCommand"`
	Status       int64    `xml:"Status"`
	Description  string   `xml:"Description"`
	Session      string   `xml:"Params>Session"`
}

//用于解析发送请求云台控制
type PTZCtrlXML struct {
	XMLName     xml.Name `xml:"PPVSPMessage"`
	Version     float64  `xml:"Version"`
	Sequence    uint64   `xml:"Sequence"`
	CommandType string   `xml:"CommandType"`
	Method      string   `xml:"Method"`
	Command     string   `xml:"Command"`
	Params      string   `xml:"Params"`
	Channel     int      `xml:"Params>Channel"`
	PTZCmd      string   `xml:"Params>PTZCmd"`
	Action      string   `xml:"Params>Action"`
	Speed       int      `xml:"Params>Speed"`
}
