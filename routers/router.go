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
			beego.NSRouter("/login", &controllers.SystemController{}, "post:Login"),
			//登出
			beego.NSRouter("/logout", &controllers.SystemController{}, "get:Logout"),
			//获取系统信息：CPU、网络、内存等
			beego.NSRouter("/info", &controllers.SystemController{}, "get:SystemInfo"),
		),

		//设备接口
		//查询设备列表
		beego.NSRouter("/devices", &controllers.DeviceController{}, "get:DeviceList"),
		//查询指定设备下的通道记录
		beego.NSRouter("/devices/:id/channels", &controllers.DeviceController{}, "get:AppointChannelList"),

		//实时直播接口
		//开始实时直播
		beego.NSRouter("/channels/:id/stream", &controllers.StreamController{}, "get:StartStream"),
	)
	beego.AddNamespace(ns)

	//Hook接口
	nsHook := beego.NewNamespace("/index",
		beego.NSNamespace("/hook",
			beego.NSRouter("/on_stream_none_reader", &controllers.HookController{}, "post:StreamNoneReaderHook"),
			beego.NSRouter("/on_publish", &controllers.HookController{}, "post:PublishHook"),
			beego.NSRouter("/on_record_mp4", &controllers.HookController{}, "post:RecordMP4Hook"),
			beego.NSRouter("/on_http_access", &controllers.HookController{}, "post:HTTPAccessHook"),
			beego.NSRouter("/on_stream_not_found", &controllers.HookController{}, "post:StreamNotFound"),
		),
	)

	beego.AddNamespace(nsHook)
}
