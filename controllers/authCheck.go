package controllers

import (
	"github.com/astaxie/beego"
	"github.com/patrickmn/go-cache"
	"github.com/tsingeye/FreeEhome/config"
)

type AuthCheckController struct {
	beego.Controller
}

func (a *AuthCheckController) Prepare() {
	controllerName, actionName := a.GetControllerAndAction()
	//fmt.Println("controllerName: ", controllerName, " actionName: ", actionName)
	//非登录时需进行登录权限验证
	if controllerName != "SystemController" || actionName != "Login" {
		a.auth()
	}
}

func (a *AuthCheckController) auth() {
	token := a.GetString("token")
	value, ok := config.AuthCheck.Get(token)
	if ok {
		//重置authCode过期时间
		config.AuthCheck.Set(token, value.(string), cache.DefaultExpiration)
	} else {
		a.Data["json"] = map[string]interface{}{
			"errCode": config.FreeEHomeUnauthorized,
			"errMsg":  config.FreeEHomeCodeMap[config.FreeEHomeUnauthorized],
		}
		a.ServeJSON()
	}
}
