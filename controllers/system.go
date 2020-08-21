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
 * @api {get} /api/v1/system/login 登录
 * @apiVersion 1.0.0
 * @apiGroup system
 * @apiName Login
 * @apiParam {String} username 用户名
 * @apiParam {String} password 密码
 * @apiSuccessExample  {json} Response-Example
 * {
 *   "errCode": 200,
 *   "errMsg": "Success OK",
 *   "authCode": "188B7DF06C77FDBE69EB25BFE946D33E" //授权码，后续所有接口需验证此授权码
 * }
 */
func (s *SystemController) Login() {
	username := s.GetString("username")
	password := s.GetString("password")

	s.Data["json"] = models.Login(username, password, s.Ctx.Request.RemoteAddr)
	s.ServeJSON()
}

/**
 * @api {get} /api/v1/system/logout 登出
 * @apiVersion 1.0.0
 * @apiGroup system
 * @apiName Logout
 * @apiParam {String} authCode 授权码
 * @apiSuccessExample  {json} Response-Example
 * {
 *   "errCode": 200,
 *   "errMsg": "Success OK"
 * }
 */
func (s *SystemController) Logout() {
	authCode := s.GetString("authCode")

	s.Data["json"] = models.Logout(authCode)
	s.ServeJSON()
}

/**
 * @api {get} /api/v1/system/info 获取系统信息
 * @apiVersion 1.0.0
 * @apiGroup system
 * @apiName Info
 * @apiParam {String} authCode 授权码
 * @apiSuccessExample  {json} Response-Example
 * {
 *   "errCode": 200,
 *   "errMsg": "Success OK",
 *   "authCode": "188B7DF06C77FDBE69EB25BFE946D33E",
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
func (s *SystemController) Info() {
	authCode := s.GetString("authCode")

	s.Data["json"] = models.SystemInfo(authCode)
	s.ServeJSON()
}
