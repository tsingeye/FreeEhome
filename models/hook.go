package models

import "github.com/tsingeye/FreeEhome/service/udp"

func StreamNotFound(deviceID, channelID string) {
	udp.StreamNotFound(deviceID, channelID)
}
