package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/tsingeye/FreeEhome/config"
	"github.com/tsingeye/FreeEhome/models"
	"github.com/tsingeye/FreeEhome/tools"
	"github.com/tsingeye/FreeEhome/tools/logs"
)

type HookController struct {
	beego.Controller
}

func (h *HookController) StopHook() {
	//fmt.Printf("%s on_stream_none_reader body: %s\n", time.Now().Format("2006-01-02 15:04:05"), string(h.Ctx.Input.RequestBody))
	logs.BeeLogger.Info("on_stream_none_reader body: %s", string(h.Ctx.Input.RequestBody))
}

//on_publish body: {
//	"app" : "rtp",
//	"id" : "1C144B0A",
//	"ip" : "192.168.1.101",
//	"params" : "",
//	"port" : 36674,
//	"schema" : "rtp",
//	"stream" : "1C144B0A",
//	"vhost" : "__defaultVhost__"
//}
func (h *HookController) PublishHook() {
	//fmt.Printf("%s on_publish body: %s\n", time.Now().Format("2006-01-02 15:04:05"), string(h.Ctx.Input.RequestBody))
	logs.BeeLogger.Info("on_publish body: %s", string(h.Ctx.Input.RequestBody))

	replyData := map[string]interface{}{
		"code":       0,
		"enableHls":  true,
		"enableMP4":  false,
		"enableRtxp": true,
		"msg":        "success",
	}

	defer func() {
		h.Data["json"] = replyData
		h.ServeJSON()
	}()

	//解析接收的内容
	data := struct {
		App    string `json:"app"` //流应用名
		Stream string `json:"stream"`
	}{}

	err := json.Unmarshal(h.Ctx.Input.RequestBody, &data)
	if err != nil {
		logs.BeeLogger.Error("analysis hook on_publish error: %s", err)
		return
	} else {
		if data.App == "download" {
			replyData["enableMP4"] = true
		}
	}

	bigInt := tools.HexToBigInt(data.Stream)
	if bigInt != nil {
		session := fmt.Sprintf("%s", bigInt)
		config.HookSession.Set(session, data.Stream, -1)
	}
}

func (h *HookController) RecordMP4Hook() {
	//fmt.Printf("%s on_record_mp4 body: %s\n", time.Now().Format("2006-01-02 15:04:05"), string(h.Ctx.Input.RequestBody))
	logs.BeeLogger.Info("on_record_mp4 body: %s", string(h.Ctx.Input.RequestBody))
}

func (h *HookController) HTTPAccessHook() {
	//fmt.Printf("%s on_http_access body: %s\n", time.Now().Format("2006-01-02 15:04:05"), string(h.Ctx.Input.RequestBody))
	logs.BeeLogger.Info("on_http_access body: %s", string(h.Ctx.Input.RequestBody))
}

//on_stream_not_found body: {
//	"app" : "rtp",
//	"id" : "140404024416544",
//	"ip" : "192.168.1.169",
//	"params" : "",
//	"port" : 60888,
//	"schema" : "rtmp",
//	"stream" : "76F94898",
//	"vhost" : "__defaultVhost__"
//}
//收到此hook指令表示流媒体重启等操作，此时内存中存在的sessionURL无效，需清除
func (h *HookController) StreamNotFound() {
	//fmt.Printf("%s on_stream_not_found body: %s\n", time.Now().Format("2006-01-02 15:04:05"), string(h.Ctx.Input.RequestBody))
	logs.BeeLogger.Info("on_stream_not_found body: %s", string(h.Ctx.Input.RequestBody))
	replyData := map[string]interface{}{
		"code": 0,
		"msg":  "success",
	}

	defer func() {
		h.Data["json"] = replyData
		h.ServeJSON()
	}()
	//解析接收的内容
	data := struct {
		Stream string `json:"stream"`
	}{}

	err := json.Unmarshal(h.Ctx.Input.RequestBody, &data)
	if err != nil {
		logs.BeeLogger.Error("analysis hook on_publish error: %s", err)
		return
	}

	bigInt := tools.HexToBigInt(data.Stream)
	if bigInt != nil {
		deviceID, channelID := config.StreamSession.Filter(fmt.Sprintf("%s", bigInt))
		config.StreamSession.Delete(channelID)

		models.StreamNotFound(deviceID, channelID)
	}
}
