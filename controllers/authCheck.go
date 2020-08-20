package controllers

import (
	"github.com/astaxie/beego"
	"github.com/patrickmn/go-cache"
	"github.com/tsingeye/FreeEhome/config"
	"github.com/tsingeye/FreeEhome/tools"
	"strings"
)

type AuthCheckController struct {
	beego.Controller
}

func (a *AuthCheckController) Prepare() {
	controllerName, actionName := a.GetControllerAndAction()
	//fmt.Println("controllerName: ", controllerName, " actionName: ", actionName)
	controllerAction := tools.StringsJoin(strings.ToLower(controllerName[0:len(controllerName)-10]), "/", strings.ToLower(actionName))
	//非登录时需进行登录权限验证
	if controllerAction != "system/login" {
		a.authCheck(controllerName, actionName)
	}
}

func (a *AuthCheckController) authCheck(controllerName, actionName string) {
	authCode := a.GetString("authCode")
	authClient, ok := config.AuthCheck.Get(authCode)
	if ok {
		//重置authCode过期时间
		config.AuthCheck.Set(authCode, authClient, cache.DefaultExpiration)
	} else {
		a.Data["json"] = map[string]interface{}{
			"errCode": config.FreeEHomeUnauthorized,
			"errMsg":  config.FreeEHomeCodeMap[config.FreeEHomeUnauthorized],
		}
	}
}
