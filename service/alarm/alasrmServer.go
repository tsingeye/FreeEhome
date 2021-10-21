package alarm

import (
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/tinyhubs/tinydom"
	"github.com/tsingeye/FreeEhome/config"
	"github.com/tsingeye/FreeEhome/tools"
	"github.com/tsingeye/FreeEhome/tools/logs"
)

//启动报警UDP服务器
func ListenAlarmUDPServer() {
	//创建udp地址
	udpAddr, err := net.ResolveUDPAddr("udp", config.AlarmAddr)
	if err != nil {
		fmt.Printf("%s ListenAlarmUDPServer() ---> net.ResolveUDPAddr() error:%s\n", time.Now().Format(config.TimeLayout), err)
		logs.PanicLogger.Panicln(fmt.Sprintf("ListenAlarmUDPServer() ---> net.ResolveUDPAddr() error:%s", err))
	}

	//监听udp地址
	udpConn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		fmt.Printf("%s ListenAlarmUDPServer() ---> net.ListenUDP() error:%s\n", time.Now().Format(config.TimeLayout), err)
		logs.PanicLogger.Panicln(fmt.Sprintf("ListenAlarmUDPServer() ---> net.ListenUDP() error:%s", err))
	}

	fmt.Printf("%s ListenAlarmUDPServer() ---> start EHome's alarmUDPServer successful!!!\n", time.Now().Format(config.TimeLayout))
	logs.BeeLogger.Emergency("ListenAlarmUDPServer() ---> start EHome's alarmUDPServer successful!!!")

	//延迟关闭监听
	defer udpConn.Close()
	for {
		buf := make([]byte, 2048)
		//阻塞获取数据
		n, _, err := udpConn.ReadFromUDP(buf)
		if err != nil {
			logs.BeeLogger.Error("ListenAlarmUDPServer() ---> udpConn.ReadFromUDP() error:%s", err)
			return
		}

		logs.BeeLogger.Debug("ListenAlarmUDPServer() ---> receive raw data:%s", buf[:n])

		go eHomeAlarmProcess(buf[:n])
	}
}

//报警信息转发
func eHomeAlarmProcess(requestBody []byte) {
	utfBody, err := tools.GB2312ToUTF8(requestBody)
	if err != nil {
		return
	}
	doc, err := tinydom.LoadDocument(strings.NewReader(string(utfBody)))
	if err != nil {
		fmt.Println("load xml fail.")
		return
	}
	//TODO 根据业务处理报文
	fmt.Println(doc)
}
