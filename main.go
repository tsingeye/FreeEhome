package main

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/plugins/cors"
	figure "github.com/common-nighthawk/go-figure"
	"github.com/kqbi/service"
	_ "github.com/tsingeye/FreeEhome/routers"
	"github.com/tsingeye/FreeEhome/service/udp"
	"github.com/tsingeye/FreeEhome/tools/logs"
	"github.com/tsingeye/FreeEhome/tools/sqlDB"
	"os"
	"time"
)

type program struct {
	exit chan bool
}

func (p *program) Start(s service.Service) error {
	go p.run()
	return nil
}

func (p *program) run() {
	figure.NewFigure("FreeEhome", "", true).Print()
	logs.BeeLogger.Info("FreeEHome Service Start!!!")
	fmt.Printf("%s FreeEHome Service Start!!!\n", time.Now().Format("2006-01-02 15:04:05"))
	//初始化数据库
	sqlDB.InitDB()
	//开启UDP服务
	go udp.ListenUDPServer()
	//实现服务端允许跨域
	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
		AllowCredentials: true,
	}))
	beego.Run()
	return
}

func (p *program) Stop(s service.Service) error {
	logs.BeeLogger.Info("FreeEHome Service Stop!!!")
	fmt.Printf("%s FreeEHome Service Stop!!!\n", time.Now().Format("2006-01-02 15:04:05"))
	close(p.exit)
	if service.Interactive() {
		os.Exit(0)
	}
	return nil
}

func main() {
	svcConfig := &service.Config{
		Name:        "FreeEHome", //服务显示名称
		DisplayName: "FreeEHome", //服务名称
		Description: "FreeEHome", //服务描述
	}

	prg := &program{
		exit: make(chan bool),
	}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		logs.PanicLogger.Fatalln("service.New() error: ", err)
	}

	if len(os.Args) > 1 {
		//install, uninstall, start, stop 的另一种实现方式
		err = service.Control(s, os.Args[1])
		if err != nil {
			logs.PanicLogger.Fatalln(err)
		}
		return
	}

	err = s.Run()
	if err != nil {
		logs.PanicLogger.Fatalln(err)
	}
}
