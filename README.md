# 海康ehome开源服务
# 简介
    EHOME协议是设备和服务器通信的一种推模式协议，适用于支持EHOME协议的网络摄像机、网络球机、DVR、NVR、车载DVR、车载取证系统、单兵、报警主机等设备。
    海康设备可以基于ehome协议来主动注册云端，区别于onvif只能在局域网内使用的限制。
    本服务软件基于海康私有协议ehome v2.x版本，力争打造一个开源安防基础产品。
# 功能
- [X] 实时预览
- [ ] 远程回放
- [X] 报警监听
- [ ] 语音对讲

# 架构
- 系统基于beego框架开发，提供RESTful接口
- CMS信令由海康ehome协议而来，基于UDP+XML进行通信
- SMS基于[ZLMediaKit](https://github.com/xia-chu/ZLMediaKit),做了二次修改，参见[MediaServer](https://github.com/kqbi/ZLMediaKit.git)
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
  
# 使用说明
## 下载程序
[v1.0](https://github.com/tsingeye/FreeEhome/releases)
## 修改CMS配置
- 进入FreeEhomeCMS=>conf文件夹
- 根据app.conf文件中注释，按实际情况修改，如下：
```ini
#ehomeCMS服务地址
udpAddr = "192.168.1.72:7660"
#流媒体SMS地址
streamStartIP = "192.168.1.72"
streamStartPort = 10000
#直播等待超时时间，默认3秒
waitStreamSessionTime = 3
#直播关闭超时时间，默认3秒
waitHookSessionTime = 3

#流媒体分发服务配置，注意必须使用英文状态下;进行分割
streamIP = "192.168.1.72;192.168.1.72"
#hls端口
hlsPort = "10080;10080"
#rtmp端口
rtmpPort = "1935;1935"
#rtsp端口
rtspPort = "10554;10554"
```
`PS.以上涉及的IP及端口按照实际配置配合SMS进行相应修改`
## 修改SMS配置
- 进入FreeEhomeSMS文件夹
- 找到config.ini文件，可选择性修改，其中流媒体分发端口需与CMS中配置一致。
- 其中【hook】部分，如果CMS和SMS在同一台机器上，可不用修复，否则这里修改为CMS的实际地址。
- 【rtp_proxy】部分配置即为海康ehome协议收流地址，需与CMS保持一致  
参考：
```ini
[rtp_proxy]
checkSource=1
dumpDir=
port=10000
timeoutSec=15
[hook]
admin_params=secret=035c73f7-bb6b-4889-a715-d9eb2d1925cc
enable=1
on_flow_report=
on_http_access=
on_play=
on_publish=http://127.0.0.1:8080/index/hook/on_publish
on_record_mp4=
on_rtsp_auth=
on_rtsp_realm=
on_server_started=
on_shell_login=
on_stream_changed=
on_stream_none_reader=http://127.0.0.1:8080/index/hook/on_stream_none_reader
on_stream_not_found=http://127.0.0.1:8080/index/hook/on_stream_not_found
timeoutSec=20
```
# 运行
`目前release版本只支持Windows平台`  
- 以终端方式运行:双击`FreeEhome.exe` `MediaServer.exe`即可
- 以服务方式运行:双击执行`install.bat`即可安装为系统服务；`uninstall.bat`为卸载系统服务；`MediaServer`暂不支持Windows系统服务。
# REST接口
参见apidoc
# 技术交流
QQ群： 1033175645
