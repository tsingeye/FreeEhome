package udp

import (
	"encoding/xml"
	"fmt"
	"github.com/tsingeye/FreeEhome/config"
	"github.com/tsingeye/FreeEhome/tools"
	"github.com/tsingeye/FreeEhome/tools/logs"
	"github.com/tsingeye/FreeEhome/tools/sqlDB"
	"net"
	"strings"
	"time"
)

//处理获得的系统服务器配置
func writeServerInfoToDB(deviceID string, utf8Data []byte, udpAddr *net.UDPAddr) {
	resGetServerInfoData := config.ResGetServerInfo{}
	if err := xml.Unmarshal(utf8Data, &resGetServerInfoData); err != nil {
		logs.BeeLogger.Error("writeServerInfoToDB() => xml.Unmarshal() error: %s", err)
		return
	}

	deviceList := &sqlDB.DeviceList{
		DeviceID:     deviceID,
		DeviceIP:     udpAddr.IP.String(),
		DeviceName:   "",
		SerialNumber: resGetServerInfoData.SerialNumber,
		Status:       "ON",
	}

	if sqlDB.Save(deviceList) {
		tabName := sqlDB.GetTableName(&sqlDB.DeviceList{})
		logs.BeeLogger.Info("deviceID=%s update record's SerialNumber to %s's table successfully", deviceID, tabName)
	}
}

//处理获得的设备工作状态信息
func writeGetDeviceWorkStatusToDB(deviceID string, utf8Data []byte, udpAddr *net.UDPAddr) {
	resGetDeviceWorkStatusData := config.ResGetDeviceWorkStatus{}
	if err := xml.Unmarshal(utf8Data, &resGetDeviceWorkStatusData); err != nil {
		logs.BeeLogger.Error("writeGetDeviceWorkStatusToDB() => xml.Unmarshal() error: %s", err)
		return
	}

	//先更新设备通道列表，将Status置为OFF，然后根据推送的设备通道工作状态重新插入或更新通道列表
	sqlDB.UpdateTable(fmt.Sprintf(`UPDATE %s SET Status = 'OFF', UpdatedAt= '%v' WHERE DeviceID= '%s'`, sqlDB.GetTableName(&sqlDB.ChannelList{}), time.Now(), deviceID))

	for _, ch := range resGetDeviceWorkStatusData.CHStatus.CHList {
		id := strings.Split(ch, "-")[0]
		channelList := &sqlDB.ChannelList{
			ChannelID:   tools.StringsJoin(deviceID, "_", id),
			ChannelName: tools.StringsJoin("Camera", id),
			DeviceID:    deviceID,
			Status:      "ON",
		}
		if sqlDB.Save(channelList) {
			tabName := sqlDB.GetTableName(&sqlDB.ChannelList{})
			logs.BeeLogger.Info("deviceID=%s saved record to %s's table successfully", deviceID, tabName)
		}
	}
}
