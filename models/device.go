package models

import (
	"github.com/tsingeye/FreeEhome/config"
	"github.com/tsingeye/FreeEhome/tools/sqlDB"
)

//分页查询设备列表：当page=limit=0时查询整个设备列表
func DeviceList(authCode string, page, limit uint64, status string) (replyData map[string]interface{}) {
	var deviceList []sqlDB.DeviceList
	var totalCount int64
	if page == 0 && limit == 0 {
		//查询整个设备列表，此时无需进行分页查询
		sqlDB.Find(&deviceList, &sqlDB.DeviceList{})
		totalCount = int64(len(deviceList))
	} else {
		//否则进行分页查询操作
		switch status {
		case "ON", "OFF":
			sqlDB.Limit(&deviceList, (page-1)*limit, limit, "Status = ?", status)
		default:
			sqlDB.Limit(&deviceList, (page-1)*limit, limit)
		}

		totalCount = sqlDB.Count(&sqlDB.DeviceList{}, "")
	}

	replyData = map[string]interface{}{
		"errCode":    config.FreeEHomeSuccessOK,
		"errMsg":     config.FreeEHomeCodeMap[config.FreeEHomeSuccessOK],
		"authCode":   authCode,
		"totalCount": totalCount,
		"deviceList": deviceList,
	}
	return
}

//分页查询设备通道列表：当page=limit=0时查询整个设备通道列表
func ChannelList(authCode, deviceID string, page, limit uint64, status string) (replyData map[string]interface{}) {
	var channelList []sqlDB.ChannelList
	var totalCount int64
	if page == 0 && limit == 0 {
		//查询整个设备通道列表，此时无需进行分页查询
		sqlDB.Find(&channelList, &sqlDB.DeviceList{})
		totalCount = int64(len(channelList))
	} else {
		//否则进行分页查询操作
		switch status {
		case "ON", "OFF":
			sqlDB.Limit(&channelList, (page-1)*limit, limit, "DeviceID = ? AND Status = ?", deviceID, status)
		default:
			sqlDB.Limit(&channelList, (page-1)*limit, limit, "DeviceID = ?", deviceID)
		}

		totalCount = sqlDB.Count(&sqlDB.DeviceList{}, "")
	}

	replyData = map[string]interface{}{
		"errCode":     config.FreeEHomeSuccessOK,
		"errMsg":      config.FreeEHomeCodeMap[config.FreeEHomeSuccessOK],
		"authCode":    authCode,
		"totalCount":  totalCount,
		"channelList": channelList,
	}
	return
}
