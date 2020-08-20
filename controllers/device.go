package controllers

import (
	"github.com/tsingeye/FreeEhome/config"
	"github.com/tsingeye/FreeEhome/models"
	"github.com/tsingeye/FreeEhome/tools/logs"
	"strings"
)

type DeviceController struct {
	AuthCheckController
}

//分页查询设备列表：当page=limit=0时查询整个设备列表
func (d *DeviceController) DeviceList() {
	authCode := d.GetString("authCode")
	replyData := map[string]interface{}{
		"errCode":  config.FreeEHomeParameterError,
		"errMsg":   config.FreeEHomeCodeMap[config.FreeEHomeParameterError],
		"authCode": authCode,
	}

	defer func() {
		d.Data["json"] = replyData
		d.ServeJSON()
	}()

	//页码编号，分页时应从1开始
	page, err := d.GetUint64("page")
	if err != nil {
		logs.BeeLogger.Error("DeviceList() --->authCode=%s, parameter <page> type error: %s", authCode, err)
		return
	}

	//分页大小
	limit, err := d.GetUint64("limit")
	if err != nil {
		logs.BeeLogger.Error("DeviceList() --->authCode=%s, parameter <limit> type error: %s", authCode, err)
		return
	}

	if page == 0 && limit != 0 {
		return
	}

	replyData = models.DeviceList(authCode, page, limit, strings.ToUpper(d.GetString("status")))
	return
}

//查询设备通道列表
func (d *DeviceController) ChannelList() {
	authCode := d.GetString("authCode")
	replyData := map[string]interface{}{
		"errCode":  config.FreeEHomeParameterError,
		"errMsg":   config.FreeEHomeCodeMap[config.FreeEHomeParameterError],
		"authCode": authCode,
	}

	defer func() {
		d.Data["json"] = replyData
		d.ServeJSON()
	}()

	//页码编号，分页时应从1开始
	page, err := d.GetUint64("page")
	if err != nil {
		logs.BeeLogger.Error("ChannelList() --->authCode=%s, parameter <page> type error: %s", authCode, err)
		return
	}

	//分页大小
	limit, err := d.GetUint64("limit")
	if err != nil {
		logs.BeeLogger.Error("ChannelList() --->authCode=%s, parameter <limit> type error: %s", authCode, err)
		return
	}

	if page == 0 && limit != 0 {
		return
	}

	replyData = models.ChannelList(authCode, d.GetString("deviceID"), page, limit, strings.ToUpper(d.GetString("status")))
	return
}
