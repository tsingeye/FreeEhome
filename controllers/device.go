package controllers

import (
	"github.com/tsingeye/FreeEhome/models"
	"strings"
)

/**
 * @apiDefine device 设备信息接口
 */

type DeviceController struct {
	AuthCheckController
}

/**
 * @api {get} /api/v1/devices 设备列表
 * @apiVersion 1.0.0
 * @apiGroup device
 * @apiName DeviceList
 * @apiParam {String} token 授权码
 * @apiParam {Number} [page] 页码，分页时默认从1开始
 * @apiParam {Number} [limit] 分页大小，默认100
 * @apiParam {String} [status] 按设备状态查询，在线：ON；离线：OFF，状态值不区分大小写，非二者则默认查询所有记录
 * @apiParam {bool} [noPage] 是否不分页，true：不分页；false：分页。布尔类型不区分大小写，默认分页
 * @apiSuccessExample  {json} Response-Example
 * {
 *   "errCode": 200,
 *   "errMsg": "Success OK",
 *   "totalCount": 100, //符合status状态的设备总数
 *   "deviceList": [
 *     {
 *       "deviceID": "ys666", //设备ID
 *       "deviceIP": "192.168.1.169", //设备IP
 *       "deviceName": "ys", //设备名
 *       "serialNumber": "666666", //设备序列号
 *       "status": "ON" //设备状态：ON-在线；OFF-离线
 *       "createdAt": "2020-10-20 10-20-10", //创建时间
 *       "updatedAt": "2020-10-20 10-20-10" //更新时间
 *     }
 *    ]
 * }
 */
func (d *DeviceController) DeviceList() {
	//页码，可选，分页时默认从1开始
	page, err := d.GetInt("page")
	if err != nil || page <= 0 {
		page = 1
	}

	//分页大小，可选，默认100
	limit, err := d.GetInt("limit")
	if err != nil || limit <= 0 {
		limit = 100
	}

	//按设备状态查询，在线：ON；离线：OFF，默认查询所有设备列表记录
	status := strings.ToUpper(d.GetString("status"))

	//是否不分页，默认分页
	noPage, err := d.GetBool("noPage")
	if err != nil {
		noPage = false
	}

	d.Data["json"] = models.DeviceList(page, limit, status, noPage)
	d.ServeJSON()
}

/**
 * @api {get} /api/v1/devices/:id/channels 查询指定设备下的通道列表
 * @apiVersion 1.0.0
 * @apiGroup device
 * @apiName AppointChannelList
 * @apiDescription 注释：:id参数是deviceID
 * @apiParam {String} token 授权码
 * @apiParam {Number} [page] 页码，分页时默认从1开始
 * @apiParam {Number} [limit] 分页大小，默认为100
 * @apiParam {String} [status] 按通道状态查询，在线：ON；离线：OFF，状态值不区分大小写，非二者则默认查询所有记录
 * @apiParam {bool} [noPage] 是否不分页，true：不分页；false：分页。布尔类型不区分大小写，默认分页
 * @apiSuccessExample  {json} Response-Example
 * {
 *   "errCode": 200,
 *   "errMsg": "Success OK",
 *   "totalCount": 100 //符合status状态的通道总数
 *   "channelList": [
 *     {
 *       "channelID": "ys666_123", //通道ID
 *       "channelName": "Camera123", //通道名
 *       "deviceID": "ys666", //设备ID
 *       "status": "ON" //设备状态：ON-在线；OFF-离线
 *       "createdAt": "2020-10-20 10-20-10", //创建时间
 *       "updatedAt": "2020-10-20 10-20-10" //更新时间
 *     }
 *   ]
 * }
 */
func (d *DeviceController) AppointChannelList() {
	//设备编号
	deviceID := d.Ctx.Input.Param(":id")

	//页码，可选，分页时默认从1开始
	page, err := d.GetInt("page")
	if err != nil || page <= 0 {
		page = 1
	}

	//分页大小，可选，默认100
	limit, err := d.GetInt("limit")
	if err != nil || limit <= 0 {
		limit = 100
	}

	//按通道状态查询，在线：ON；离线：OFF，默认查询所有记录
	status := strings.ToUpper(d.GetString("status"))

	//是否不分页，默认分页
	noPage, err := d.GetBool("noPage")
	if err != nil {
		noPage = false
	}

	d.Data["json"] = models.AppointChannelList(deviceID, page, limit, status, noPage)
	d.ServeJSON()
}
