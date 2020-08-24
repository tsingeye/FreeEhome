package controllers

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/tsingeye/FreeEhome/tools/logs"
	"time"
)

type HookController struct {
	beego.Controller
}

func (h *HookController) StopHook() {
	fmt.Printf("%s on_stream_none_reader body: %s\n", time.Now().Format("2006-01-02 15:04:05"), string(h.Ctx.Input.RequestBody))
	logs.BeeLogger.Info("on_stream_none_reader body: %s", string(h.Ctx.Input.RequestBody))
}

func (h *HookController) PublishHook() {
	fmt.Printf("%s on_publish body: %s\n", time.Now().Format("2006-01-02 15:04:05"), string(h.Ctx.Input.RequestBody))
	logs.BeeLogger.Info("on_publish body: %s", string(h.Ctx.Input.RequestBody))
}

func (h *HookController) RecordMP4Hook() {
	fmt.Printf("%s on_record_mp4 body: %s\n", time.Now().Format("2006-01-02 15:04:05"), string(h.Ctx.Input.RequestBody))
	logs.BeeLogger.Info("on_record_mp4 body: %s", string(h.Ctx.Input.RequestBody))
}

func (h *HookController) HTTPAccessHook() {
	fmt.Printf("%s on_http_access body: %s\n", time.Now().Format("2006-01-02 15:04:05"), string(h.Ctx.Input.RequestBody))
	logs.BeeLogger.Info("on_http_access body: %s", string(h.Ctx.Input.RequestBody))
}
