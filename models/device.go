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

//统计设备通道列表符合status条件的数量
func channelCount(deviceID string, status string) (count int64) {
	switch deviceID {
	case "":
		//查询通道列表中符合条件的所有记录
		switch status {
		case "ON", "OFF":
			count = sqlDB.Count(&sqlDB.ChannelList{}, "Status = ?", status)
		default:
			count = sqlDB.Count(&sqlDB.ChannelList{}, "")
		}
	default:
		//查询通道列表中指定deviceID的符合条件的记录
		switch status {
		case "ON", "OFF":
			count = sqlDB.Count(&sqlDB.ChannelList{}, "DeviceID = ? AND Status = ?", deviceID, status)
		default:
			count = sqlDB.Count(&sqlDB.ChannelList{}, "DeviceID = ?", deviceID)
		}
	}

	return
}

//查询设备通道列表
func ChannelList(deviceID string, page, limit int, status string, noPage bool) (replyData map[string]interface{}) {
	var channelList []sqlDB.ChannelList
	switch deviceID {
	case "":
		//查询通道列表中符合条件的所有记录
		switch noPage {
		case true:
			//不使用分页查询
			switch status {
			case "ON", "OFF":
				sqlDB.Find(&channelList, sqlDB.GetTableName(&sqlDB.ChannelList{}), "Status = ?", status)
			default:
				sqlDB.Find(&channelList, sqlDB.GetTableName(&sqlDB.ChannelList{}))
			}
		case false:
			//分页查询
			offset := (page - 1) * limit
			switch status {
			case "ON", "OFF":
				sqlDB.Limit(&channelList, sqlDB.GetTableName(&sqlDB.ChannelList{}), limit, offset, "Status = ?", status)
			default:
				sqlDB.Limit(&channelList, sqlDB.GetTableName(&sqlDB.ChannelList{}), limit, offset)
			}
		}
	default:
		//查询通道列表中指定deviceID的符合条件的记录
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
	}
	
	if len(channelList) == 0 {
		channelList = make([]sqlDB.ChannelList, 0)
	}

	replyData = map[string]interface{}{
		"errCode":     config.FreeEHomeSuccessOK,
		"errMsg":      config.FreeEHomeCodeMap[config.FreeEHomeSuccessOK],
		"totalCount":  channelCount(deviceID, status),
		"channelList": channelList,
	}
	return
}
