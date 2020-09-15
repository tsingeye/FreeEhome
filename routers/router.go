package routers

import (
	"github.com/astaxie/beego"
	"github.com/tsingeye/FreeEhome/controllers"
)

func init() {
	ns := beego.NewNamespace("/api/v1",
		//系统接口
		beego.NSNamespace("/system",
			//登录
			beego.NSRouter("/login", &controllers.SystemController{}, "get:Login"),
			//登出
			beego.NSRouter("/logout", &controllers.SystemController{}, "get:Logout"),
			//获取系统信息：CPU、网络、内存等
			beego.NSRouter("/info", &controllers.SystemController{}, "get:Info"),
		),
		//设备接口
		beego.NSNamespace("/device",
			//查询设备列表
			beego.NSRouter("/list", &controllers.DeviceController{}, "get:DeviceList"),
			//查询通道列表
			beego.NSRouter("/channelList", &controllers.DeviceController{}, "get:ChannelList"),
		),
		//实时直播接口
		beego.NSNamespace("/stream",
			//开始实时直播
			beego.NSRouter("/start", &controllers.StreamController{}, "get:Start"),
			//关闭实时直播
			beego.NSRouter("/stop", &controllers.StreamController{}, "get:Stop"),
		),
	)
	beego.AddNamespace(ns)

	//Hook接口
	nsHook := beego.NewNamespace("/index",
		beego.NSNamespace("/hook",
			beego.NSRouter("/on_stream_none_reader", &controllers.HookController{}, "post:StopHook"),
			beego.NSRouter("/on_publish", &controllers.HookController{}, "post:PublishHook"),
			beego.NSRouter("/on_record_mp4", &controllers.HookController{}, "post:RecordMP4Hook"),
			beego.NSRouter("/on_http_access", &controllers.HookController{}, "post:HTTPAccessHook"),
			beego.NSRouter("/on_stream_not_found", &controllers.HookController{}, "post:StreamNotFound"),
		),
	)

	beego.AddNamespace(nsHook)
}
