package controllers

import (
	"github.com/tsingeye/FreeEhome/models"
)

type SystemController struct {
	AuthCheckController
}

//登录
func (s *SystemController) Login() {
	username := s.GetString("username")
	password := s.GetString("password")

	s.Data["json"] = models.Login(username, password, s.Ctx.Request.RemoteAddr)
	s.ServeJSON()
}

//登出
func (s *SystemController) Logout() {
	authCode := s.GetString("authCode")

	s.Data["json"] = models.Logout(authCode)
	s.ServeJSON()
}

//获取系统信息：CPU、网络、内存等
func (s *SystemController) Info() {
	authCode := s.GetString("authCode")

	s.Data["json"] = models.SystemInfo(authCode)
	s.ServeJSON()
}
