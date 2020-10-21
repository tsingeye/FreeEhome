package models

import "github.com/tsingeye/FreeEhome/service/udp"

func StreamNoneReaderHook(deviceID, channelID string) {
	udp.SendByeStream("on_stream_none_reader---> ", deviceID, channelID)
}

func StreamNotFound(deviceID, channelID string) {
	udp.StreamNotFound(deviceID, channelID)
}
