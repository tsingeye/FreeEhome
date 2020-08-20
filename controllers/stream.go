package controllers

import "github.com/tsingeye/FreeEhome/models"

type StreamController struct {
	AuthCheckController
}

func (s *StreamController) Start() {
	authCode := s.GetString("authCode")
	deviceID := s.GetString("deviceID")
	channelID := s.GetString("channelID")

	s.Data["json"] = models.StartStream(authCode, deviceID, channelID)
	s.ServeJSON()
}

func (s *StreamController) Stop() {
	authCode := s.GetString("authCode")
	deviceID := s.GetString("deviceID")
	channelID := s.GetString("channelID")

	s.Data["json"] = models.StopStream(authCode, deviceID, channelID)
	s.ServeJSON()
}
