package models

import (
	"github.com/tsingeye/FreeEhome/config"
	"github.com/tsingeye/FreeEhome/tools/sqlDB"
)

//统计设备列表中符合status条件的总数
func deviceCount(status string) (count int64) {
	switch status {
	case "ON", "OFF":
		count = sqlDB.Count(&sqlDB.DeviceList{}, "Status = ?", status)
	default:
		count = sqlDB.Count(&sqlDB.DeviceList{}, "")
	}
	return
}

//查询设备列表
func DeviceList(page, limit int, status string, noPage bool) (replyData map[string]interface{}) {
	var deviceList []sqlDB.DeviceList

	switch noPage {
	case true:
		//不使用分页查询
		switch status {
		case "ON", "OFF":
			sqlDB.Find(&deviceList, sqlDB.GetTableName(&sqlDB.DeviceList{}), "Status = ?", status)
		default:
			sqlDB.Find(&deviceList, sqlDB.GetTableName(&sqlDB.DeviceList{}))
		}
	case false:
		//分页查询
		offset := (page - 1) * limit
		switch status {
		case "ON", "OFF":
			sqlDB.Limit(&deviceList, sqlDB.GetTableName(&sqlDB.DeviceList{}), limit, offset, "Status = ?", status)
		default:
			sqlDB.Limit(&deviceList, sqlDB.GetTableName(&sqlDB.DeviceList{}), limit, offset)
		}
	}

	if len(deviceList) == 0 {
		deviceList = make([]sqlDB.DeviceList, 0)
	}

	replyData = map[string]interface{}{
		"errCode":    config.FreeEHomeSuccessOK,
		"errMsg":     config.FreeEHomeCodeMap[config.FreeEHomeSuccessOK],
		"totalCount": deviceCount(status),
		"deviceList": deviceList,
	}
	return
}

//查询指定设备下的符合status条件的通道总数
func appointChannelCount(deviceID string, status string) (count int64) {
	switch status {
	case "ON", "OFF":
		count = sqlDB.Count(&sqlDB.ChannelList{}, "DeviceID = ? AND Status = ?", deviceID, status)
	default:
		count = sqlDB.Count(&sqlDB.ChannelList{}, "DeviceID = ?", deviceID)
	}

	return
}

//查询指定设备下的通道列表
func AppointChannelList(deviceID string, page, limit int, status string, noPage bool) (replyData map[string]interface{}) {
	var channelList []sqlDB.ChannelList

	switch noPage {
	case true:
		//不使用分页查询
		switch status {
		case "ON", "OFF":
			sqlDB.Find(&channelList, sqlDB.GetTableName(&sqlDB.ChannelList{}), "DeviceID = ? AND Status = ?", deviceID, status)
		default:
			sqlDB.Find(&channelList, sqlDB.GetTableName(&sqlDB.ChannelList{}), "DeviceID = ?", deviceID)
		}
	case false:
		//分页查询
		offset := (page - 1) * limit
		switch status {
		case "ON", "OFF":
			sqlDB.Limit(&channelList, sqlDB.GetTableName(&sqlDB.ChannelList{}), limit, offset, "DeviceID = ? AND Status = ?", deviceID, status)
		default:
			sqlDB.Limit(&channelList, sqlDB.GetTableName(&sqlDB.ChannelList{}), limit, offset, "DeviceID = ?", deviceID)
		}
	}

	if len(channelList) == 0 {
		channelList = make([]sqlDB.ChannelList, 0)
	}

	replyData = map[string]interface{}{
		"errCode":     config.FreeEHomeSuccessOK,
		"errMsg":      config.FreeEHomeCodeMap[config.FreeEHomeSuccessOK],
		"totalCount":  appointChannelCount(deviceID, status),
		"channelList": channelList,
	}
	return
}
