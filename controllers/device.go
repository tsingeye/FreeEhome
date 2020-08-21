package controllers

import (
	"github.com/tsingeye/FreeEhome/config"
	"github.com/tsingeye/FreeEhome/models"
	"github.com/tsingeye/FreeEhome/tools/logs"
	"strings"
)

/**
 * @apiDefine device 设备信息接口
 */

type DeviceController struct {
	AuthCheckController
}

/**
 * @api {get} /api/v1/device/list 分页查询设备列表
 * @apiVersion 1.0.0
 * @apiGroup device
 * @apiName DeviceList
 * @apiParam {String} authCode 授权码
 * @apiParam {Number} page 页码，分页时应从1开始，当page=limit=0时不分页查询所有设备信息
 * @apiParam {Number} limit 分页大小
 * @apiParam {String} [status] 设备状态，ON-在线；OFF-离线，其它则分页查询所有设备信息
 * @apiSuccessExample  {json} Response-Example
 * {
 *   "errCode": 200,
 *   "errMsg": "Success OK",
 *   "authCode": "188B7DF06C77FDBE69EB25BFE946D33E" //授权码
 *   "totalCount": 100, //设备列表总数
 *   "deviceList": [
 *     {
 *       "deviceID": "ys666", //设备ID
 *       "deviceIP": "192.168.1.169", //设备IP
 *       "deviceName": "ys", //设备名
 *       "serialNumber": "666666", //设备序列号
 *       "status": "ON" //设备状态：ON-在线；OFF-离线
 *     }
 *    ]
 * }
 */
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

/**
 * @api {get} /api/v1/device/channelList 分页查询设备通道列表
 * @apiVersion 1.0.0
 * @apiGroup device
 * @apiName ChannelList
 * @apiParam {String} authCode 授权码
 * @apiParam {String} deviceID 设备ID
 * @apiParam {Number} page 页码，分页时应从1开始，当page=limit=0时不分页查询所有设备通道信息
 * @apiParam {Number} limit 分页大小
 * @apiParam {String} [status] 设备状态，ON-在线；OFF-离线，其它则分页查询所有设备通道信息
 * @apiSuccessExample  {json} Response-Example
 * {
 *   "errCode": 200,
 *   "errMsg": "Success OK",
 *   "authCode": "188B7DF06C77FDBE69EB25BFE946D33E" //授权码
 *   "totalCount": 100 //设备通道列表总数
 *   "deviceList": [
 *     {
 *       "channelID": "ys666_123", //通道ID
 *       "channelName": "Camera123", //通道名
 *       "deviceID": "ys666", //设备ID
 *       "status": "ON" //设备状态：ON-在线；OFF-离线
 *     }
 *   ]
 * }
 */
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
