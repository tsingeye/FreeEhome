# ehome
海康ehome开源服务
# 简介
    EHOME协议是设备和服务器通信的一种推模式协议，适用于支持EHOME协议的网络摄像机、网络球机、DVR、NVR、车载DVR、车载取证系统、单兵、报警主机等设备。
    海康设备可以基于ehome协议来主动注册云端，区别于onvif只能在局域网内使用的限制。
    本服务软件基于海康私有协议ehome v2.x版本，力争打造一个开源安防基础产品。
# 功能
- 实时预览
- 远程回放
- 报警监听
- 语音对讲

# 架构
- 基于beego框架开发
# 编译步骤
- go get:
    在项目中使用go get下载安装依赖包
- go build:
    - 在main.go函数所在路径上使用go build编译。
    - Windows：执行FreeEhome.exe即可运行程序；
    - Linux：执行./FreeEhome即可运行程序。
# 打包部署，使用bee工具轻松跨平台打包部署：
- Windows：bee pack -be GOOS=windows;
- Linux：bee pack -be GOOS=linux;
- Windows如何安装成服务？
    - 解压打包的文件；
    - 双击install.bat安装成服务；
    - 双击uninstall.bat卸载服务。