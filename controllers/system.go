package controllers

import (
	"github.com/tsingeye/FreeEhome/models"
)

/**
 * @apiDefine system 系统接口
 */

type SystemController struct {
	AuthCheckController
}

/**
 * @api {post} /api/v1/system/login 登录
 * @apiGroup system
 * @apiName Login
 * @apiParamExample {json} Request-Example:
 *   {
 *     "username": "wyd",
 *     "password": "wyd666" //32位MD5加密小写数据，暂时用户名密码不进行校验
 *   }
 * @apiSuccessExample {json} Success-Response:
 *   {
 *     "errCode": 200,
 *     "errMsg": "Success OK",
 *     "token": "this is token" //后面所有接口需验证此token，用法是作为URL参数使用
 *   }
 */
func (s *SystemController) Login() {
	s.Data["json"] = models.Login(s.Ctx.Request.RemoteAddr, s.Ctx.Input.RequestBody)
	s.ServeJSON()
}

/**
 * @api {get} /api/v1/system/logout 登出
 * @apiVersion 1.0.0
 * @apiGroup system
 * @apiName Logout
 * @apiParam {String} token 授权码
 * @apiSuccessExample  {json} Response-Example
 * {
 *   "errCode": 200,
 *   "errMsg": "Success OK"
 * }
 */
func (s *SystemController) Logout() {
	token := s.GetString("token")

	s.Data["json"] = models.Logout(token)
	s.ServeJSON()
}

/**
 * @api {get} /api/v1/system/info 获取系统信息
 * @apiVersion 1.0.0
 * @apiGroup system
 * @apiName SystemInfo
 * @apiParam {String} token 授权码
 * @apiSuccessExample  {json} Response-Example
 * {
 *   "errCode": 200,
 *   "errMsg": "Success OK",
 *   "cpuUsedPercent": "6%", //CPU使用率
 *   "virtualMemory": {
 *     "total": "8079MB", //总内存
 *     "available": "2565MB", //当前可用内存
 *     "used": "5514MB", //当前已使用内存
 *     "usedPercent": "68%" //当前内存使用率
 *   },
 *   "network": {
 *     "uploadSpeed": "0KB/s", //上传速度
 *     "downloadSpeed": "0KB/s" //下载速度
 *   },
 *   "deviceInfo": {
 *     "totalCount": 0, //设备总数
 *     "onlineCount": 0 //设备在线总数
 *   },
 *   "channelInfo": {
 *     "totalCount": 0, //通道总数
 *     "onlineCount": 0 //通道在线总数
 *   }
 * }
 */
func (s *SystemController) SystemInfo() {
	s.Data["json"] = models.SystemInfo()
	s.ServeJSON()
}
