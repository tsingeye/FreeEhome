package models

import (
	"strconv"
	"strings"

	"github.com/tsingeye/FreeEhome/config"
	"github.com/tsingeye/FreeEhome/service/udp"
)

func PTZCtrl(channelID, cmd, action string, speed int) (replyData map[string]interface{}) {
	replyData = map[string]interface{}{
		"errCode": config.FreeEHomeSuccessOK,
		"errMsg":  config.FreeEHomeCodeMap[config.FreeEHomeSuccessOK],
	}

	data := strings.Split(channelID, "_")
	if len(data) != 2 {
		replyData = map[string]interface{}{
			"errCode": config.FreeEHomeParameterError,
			"errMsg":  config.FreeEHomeCodeMap[config.FreeEHomeParameterError],
		}
		return
	}
	device := data[0]
	channel, _ := strconv.Atoi(data[1])
	go udp.SendPTZCtrl(device, cmd, action, channel, speed)

	return
}
