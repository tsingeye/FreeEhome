package sqlDB

import (
	"github.com/tsingeye/FreeEhome/config"
)

//设备列表
type DeviceList struct {
	DeviceID     string            `gorm:"column:DeviceID;primary_key" json:"deviceID"` //设备ID
	DeviceIP     string            `gorm:"column:DeviceIP" json:"deviceIP"`             //来自设备注册时的UDPAddr
	DeviceName   string            `gorm:"column:DeviceName" json:"deviceName"`
	SerialNumber string            `gorm:"column:SerialNumber" json:"serialNumber"` //设备序列号
	Status       string            `gorm:"column:Status" json:"status"`             //设备在线状态，ON-在线，OFF-离线
	CreatedAt    config.TimeNormal `gorm:"column:CreatedAt" json:"-"`
	UpdatedAt    config.TimeNormal `gorm:"column:UpdatedAt" json:"-"`
}

//设置表名
func (DeviceList) TableName() string {
	return "deviceList"
}

//通道列表
type ChannelList struct {
	ChannelID   string            `gorm:"column:ChannelID;primary_key" json:"channelID"` //ChannelID=DeviceID_Number
	ChannelName string            `gorm:"column:ChannelName" json:"channelName"`
	DeviceID    string            `gorm:"column:DeviceID" json:"deviceID"`
	Status      string            `gorm:"column:Status" json:"status"` //通道在线状态，ON-在线，OFF-离线
	CreatedAt   config.TimeNormal `gorm:"column:CreatedAt" json:"-"`
	UpdatedAt   config.TimeNormal `gorm:"column:UpdatedAt" json:"-"`
}

//设置表名
func (ChannelList) TableName() string {
	return "channelList"
}
