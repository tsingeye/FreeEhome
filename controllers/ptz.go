package controllers

import (
	"github.com/tsingeye/FreeEhome/models"
)

/**
 * @apiDefine ptz 云台控制接口
 */
type PTZController struct {
	AuthCheckController
}

/**
* @api {post} /api/v1/channels/:id/ptz 云台控制
* @apiVersion 1.0.0
* @apiGroup ptz
* @apiName PTZCtrl
 * @apiDescription 注释：:id参数是channelID
 * @apiParam {String} token 授权码
 * @apiParam {String} [cmd] 方向命令：LEFT RIHGT UP DOWN，空为LEFT
 * @apiParam {String} [action] 控制动作：Start、Stop，空为Start
 * @apiParam {Number} [speed] 云台控制速度：1-7，空为4
 * @apiSuccessExample {json} Success-Response:
 *   {
 *     "errCode": 200,
 *     "errMsg": "Success OK",
 *   }
*/
func (s *PTZController) PTZCtrl() {
	//通道编号
	channelID := s.Ctx.Input.Param(":id")
	cmd := s.GetString("cmd", "LEFT")
	action := s.GetString("action", "Start")
	speed, err := s.GetInt("speed")
	if err != nil {
		speed = 4
	}

	s.Data["json"] = models.PTZCtrl(channelID, cmd, action, speed)
	s.ServeJSON()
}
