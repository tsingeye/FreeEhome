package models

import (
	"github.com/tsingeye/FreeEhome/config"
	"github.com/tsingeye/FreeEhome/service/udp"
)

func StartStream(authCode, deviceID, channelID string) (replyData map[string]interface{}) {
	errCode, streamURL := udp.SendInviteStream(authCode, deviceID, channelID)
	replyData = map[string]interface{}{
		"errCode":   errCode,
		"errMsg":    config.FreeEHomeCodeMap[errCode],
		"authCode":  authCode,
		"streamURL": streamURL,
	}
	return
}

func StopStream(authCode, deviceID, channelID string) (replyData map[string]interface{}) {
	replyData = map[string]interface{}{
		"errCode":  config.FreeEHomeSuccessOK,
		"errMsg":   config.FreeEHomeCodeMap[config.FreeEHomeSuccessOK],
		"authCode": authCode,
	}
	return
}
