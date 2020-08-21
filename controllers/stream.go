package controllers

import "github.com/tsingeye/FreeEhome/models"

/**
 * @apiDefine stream 实时直播接口
 */

type StreamController struct {
	AuthCheckController
}

/**
 * @api {get} /api/v1/stream/start 开始实时直播
 * @apiVersion 1.0.0
 * @apiGroup stream
 * @apiName Start
 * @apiParam {String} authCode 授权码
 * @apiParam {String} deviceID 设备ID
 * @apiParam {String} channelID 通道ID
 * @apiSuccessExample  {json} Response-Example
 * {
 *   "errCode": 200,
 *   "errMsg": "Success OK",
 *   "authCode": "188B7DF06C77FDBE69EB25BFE946D33E", //授权码
 *   "streamURL": "https://www.baidu.com/" //实时直播URL
 * }
 */
func (s *StreamController) Start() {
	authCode := s.GetString("authCode")
	deviceID := s.GetString("deviceID")
	channelID := s.GetString("channelID")

	s.Data["json"] = models.StartStream(authCode, deviceID, channelID)
	s.ServeJSON()
}

/**
 * @api {get} /api/v1/stream/stop 关闭实时直播
 * @apiVersion 1.0.0
 * @apiGroup stream
 * @apiName Stop
 * @apiParam {String} authCode 授权码
 * @apiParam {String} deviceID 设备ID
 * @apiParam {String} channelID 通道ID
 * @apiSuccessExample  {json} Response-Example
 * {
 *   "errCode": 200,
 *   "errMsg": "Success OK",
 *   "authCode": "188B7DF06C77FDBE69EB25BFE946D33E" //授权码
 * }
 */
func (s *StreamController) Stop() {
	authCode := s.GetString("authCode")
	deviceID := s.GetString("deviceID")
	channelID := s.GetString("channelID")

	s.Data["json"] = models.StopStream(authCode, deviceID, channelID)
	s.ServeJSON()
}
