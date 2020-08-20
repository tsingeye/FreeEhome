package udp

import (
	"fmt"
	"github.com/tsingeye/FreeEhome/config"
	"github.com/tsingeye/FreeEhome/tools"
	"github.com/tsingeye/FreeEhome/tools/logs"
	"net"
	"time"
)

func ListenUDPServer() {
	udpAddr, err := net.ResolveUDPAddr("udp", config.UDPAddr)
	if err != nil {
		fmt.Println(time.Now().Format("2006-01-02 15:04:05"), "start udpServer failed, trigger panic!")
		logs.PanicLogger.Panicln(fmt.Sprintf("net.ResolveUDPAddr() error: %s", err))
	}
	udpConn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		fmt.Println(time.Now().Format("2006-01-02 15:04:05"), "start udpServer failed, trigger panic!")
		logs.PanicLogger.Panicln(fmt.Sprintf("net.ListenUDP() error: %s", err))
	}

	logs.BeeLogger.Info("start UDPServer successful!")
	fmt.Println(time.Now().Format("2006-01-02 15:04:05"), "start udpServer successful!")

	//获得一个UDP网关
	gw = &GatewayForUDP{
		UDPClientList: make(map[string]*ClientForUDP),
		Sequence:      uint64(tools.GetRangeNum(700, 3564)),
		SequenceMap:   make(map[uint64]*SequenceInfo),
		UDPConn:       udpConn,
	}

	//分发处理接收到的数据
	gw.readFromUDP()
}
