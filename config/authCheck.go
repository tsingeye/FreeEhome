package config

import (
	"github.com/astaxie/beego"
	"github.com/patrickmn/go-cache"
	"github.com/tsingeye/FreeEhome/tools/logs"
	"time"
)

var (
	AuthCheck *cache.Cache
)

func init() {
	authCodeDefaultExpiration, err := beego.AppConfig.Int64("authCodeDefaultExpiration")
	if authCodeDefaultExpiration <= 0 || err != nil {
		authCodeDefaultExpiration = 3600
		logs.BeeLogger.Error("set default value, authCodeDefaultExpiration=%d, init authCodeDefaultExpiration error: %s", authCodeDefaultExpiration, err)
	}

	authCodeCleanupInterval, err := beego.AppConfig.Int64("authCodeCleanupInterval")
	if authCodeCleanupInterval <= 0 || err != nil {
		authCodeCleanupInterval = 360
		logs.BeeLogger.Error("set default value, authCodeCleanupInterval=%d, init authCodeCleanupInterval error: %s", authCodeCleanupInterval, err)
	}

	AuthCheck = cache.New(time.Duration(authCodeDefaultExpiration)*time.Second, time.Duration(authCodeCleanupInterval)*time.Second)
}

type AuthClient struct {
	Username     string            //登录用户名
	Password     string            //登录密码
	RemoteAddr   string            //客户端地址
	StreamURLMap map[string]string //返回实时直播获得的streamURL
	AuthCode     string            //鉴权校验码
}
