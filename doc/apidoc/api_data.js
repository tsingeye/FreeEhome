define({ "api": [
  {
    "type": "get",
    "url": "/api/v1/devices/:id/channels",
    "title": "查询指定设备下的通道列表",
    "version": "1.0.0",
    "group": "device",
    "name": "AppointChannelList",
    "description": "<p>注释：:id参数是deviceID</p>",
    "parameter": {
      "fields": {
        "Parameter": [
          {
            "group": "Parameter",
            "type": "String",
            "optional": false,
            "field": "token",
            "description": "<p>授权码</p>"
          },
          {
            "group": "Parameter",
            "type": "Number",
            "optional": true,
            "field": "page",
            "description": "<p>页码，分页时默认从1开始</p>"
          },
          {
            "group": "Parameter",
            "type": "Number",
            "optional": true,
            "field": "limit",
            "description": "<p>分页大小，默认为100</p>"
          },
          {
            "group": "Parameter",
            "type": "String",
            "optional": true,
            "field": "status",
            "description": "<p>按通道状态查询，在线：ON；离线：OFF，状态值不区分大小写，非二者则默认查询所有记录</p>"
          },
          {
            "group": "Parameter",
            "type": "bool",
            "optional": true,
            "field": "noPage",
            "description": "<p>是否不分页，true：不分页；false：分页。布尔类型不区分大小写，默认分页</p>"
          }
        ]
      }
    },
    "success": {
      "examples": [
        {
          "title": "Response-Example",
          "content": "{\n  \"errCode\": 200,\n  \"errMsg\": \"Success OK\",\n  \"totalCount\": 100 //符合status状态的通道总数\n  \"channelList\": [\n    {\n      \"channelID\": \"ys666_123\", //通道ID\n      \"channelName\": \"Camera123\", //通道名\n      \"deviceID\": \"ys666\", //设备ID\n      \"status\": \"ON\" //设备状态：ON-在线；OFF-离线\n      \"createdAt\": \"2020-10-20 10-20-10\", //创建时间\n      \"updatedAt\": \"2020-10-20 10-20-10\" //更新时间\n    }\n  ]\n}",
          "type": "json"
        }
      ]
    },
    "filename": "controllers/device.go",
    "groupTitle": "设备信息接口"
  },
  {
    "type": "get",
    "url": "/api/v1/devices",
    "title": "设备列表",
    "version": "1.0.0",
    "group": "device",
    "name": "DeviceList",
    "parameter": {
      "fields": {
        "Parameter": [
          {
            "group": "Parameter",
            "type": "String",
            "optional": false,
            "field": "token",
            "description": "<p>授权码</p>"
          },
          {
            "group": "Parameter",
            "type": "Number",
            "optional": true,
            "field": "page",
            "description": "<p>页码，分页时默认从1开始</p>"
          },
          {
            "group": "Parameter",
            "type": "Number",
            "optional": true,
            "field": "limit",
            "description": "<p>分页大小，默认100</p>"
          },
          {
            "group": "Parameter",
            "type": "String",
            "optional": true,
            "field": "status",
            "description": "<p>按设备状态查询，在线：ON；离线：OFF，状态值不区分大小写，非二者则默认查询所有记录</p>"
          },
          {
            "group": "Parameter",
            "type": "bool",
            "optional": true,
            "field": "noPage",
            "description": "<p>是否不分页，true：不分页；false：分页。布尔类型不区分大小写，默认分页</p>"
          }
        ]
      }
    },
    "success": {
      "examples": [
        {
          "title": "Response-Example",
          "content": "{\n  \"errCode\": 200,\n  \"errMsg\": \"Success OK\",\n  \"totalCount\": 100, //符合status状态的设备总数\n  \"deviceList\": [\n    {\n      \"deviceID\": \"ys666\", //设备ID\n      \"deviceIP\": \"192.168.1.169\", //设备IP\n      \"deviceName\": \"ys\", //设备名\n      \"serialNumber\": \"666666\", //设备序列号\n      \"status\": \"ON\" //设备状态：ON-在线；OFF-离线\n      \"createdAt\": \"2020-10-20 10-20-10\", //创建时间\n      \"updatedAt\": \"2020-10-20 10-20-10\" //更新时间\n    }\n   ]\n}",
          "type": "json"
        }
      ]
    },
    "filename": "controllers/device.go",
    "groupTitle": "设备信息接口"
  },
  {
    "type": "post",
    "url": "/api/v1/channels/:id/ptz",
    "title": "云台控制",
    "version": "1.0.0",
    "group": "ptz",
    "name": "PTZCtrl",
    "description": "<p>注释：:id参数是channelID</p>",
    "parameter": {
      "fields": {
        "Parameter": [
          {
            "group": "Parameter",
            "type": "String",
            "optional": false,
            "field": "token",
            "description": "<p>授权码</p>"
          },
          {
            "group": "Parameter",
            "type": "String",
            "optional": true,
            "field": "cmd",
            "description": "<p>方向命令：LEFT RIHGT UP DOWN，空为LEFT</p>"
          },
          {
            "group": "Parameter",
            "type": "String",
            "optional": true,
            "field": "action",
            "description": "<p>控制动作：Start、Stop，空为Start</p>"
          },
          {
            "group": "Parameter",
            "type": "Number",
            "optional": true,
            "field": "speed",
            "description": "<p>云台控制速度：1-7，空为4</p>"
          }
        ]
      }
    },
    "success": {
      "examples": [
        {
          "title": "Success-Response:",
          "content": "{\n  \"errCode\": 200,\n  \"errMsg\": \"Success OK\",\n}",
          "type": "json"
        }
      ]
    },
    "filename": "controllers/ptz.go",
    "groupTitle": "云台控制接口"
  },
  {
    "type": "get",
    "url": "/api/v1/channels/:id/stream",
    "title": "开始实时直播",
    "version": "1.0.0",
    "group": "stream",
    "name": "StartStream",
    "description": "<p>注释：:id参数是channelID</p>",
    "parameter": {
      "fields": {
        "Parameter": [
          {
            "group": "Parameter",
            "type": "String",
            "optional": false,
            "field": "token",
            "description": "<p>授权码</p>"
          }
        ]
      }
    },
    "success": {
      "examples": [
        {
          "title": "Response-Example",
          "content": "{\n  \"errCode\": 200,\n  \"errMsg\": \"Success OK\",\n  \"sessionURL\": {\n    \"rtmp\": \"rtmp://ip:port/rtp/xxx\",\n    \"flv\": \"http://ip:port/rtp/xxx.flv\",\n    \"rtsp\": \"rtsp://ip:port/rtp/xxx\",\n    \"hls\": \"http://ip:port/rtp/xxx/hls.m3u8\",\n  }\n}",
          "type": "json"
        }
      ]
    },
    "filename": "controllers/stream.go",
    "groupTitle": "实时直播接口"
  },
  {
    "type": "post",
    "url": "/api/v1/system/login",
    "title": "登录",
    "group": "system",
    "name": "Login",
    "parameter": {
      "examples": [
        {
          "title": "Request-Example:",
          "content": "{\n  \"username\": \"wyd\",\n  \"password\": \"wyd666\" //32位MD5加密小写数据，暂时用户名密码不进行校验\n}",
          "type": "json"
        }
      ]
    },
    "success": {
      "examples": [
        {
          "title": "Success-Response:",
          "content": "{\n  \"errCode\": 200,\n  \"errMsg\": \"Success OK\",\n  \"token\": \"this is token\" //后面所有接口需验证此token，用法是作为URL参数使用\n}",
          "type": "json"
        }
      ]
    },
    "version": "0.0.0",
    "filename": "controllers/system.go",
    "groupTitle": "系统接口"
  },
  {
    "type": "get",
    "url": "/api/v1/system/logout",
    "title": "登出",
    "version": "1.0.0",
    "group": "system",
    "name": "Logout",
    "parameter": {
      "fields": {
        "Parameter": [
          {
            "group": "Parameter",
            "type": "String",
            "optional": false,
            "field": "token",
            "description": "<p>授权码</p>"
          }
        ]
      }
    },
    "success": {
      "examples": [
        {
          "title": "Response-Example",
          "content": "{\n  \"errCode\": 200,\n  \"errMsg\": \"Success OK\"\n}",
          "type": "json"
        }
      ]
    },
    "filename": "controllers/system.go",
    "groupTitle": "系统接口"
  },
  {
    "type": "get",
    "url": "/api/v1/system/info",
    "title": "获取系统信息",
    "version": "1.0.0",
    "group": "system",
    "name": "SystemInfo",
    "parameter": {
      "fields": {
        "Parameter": [
          {
            "group": "Parameter",
            "type": "String",
            "optional": false,
            "field": "token",
            "description": "<p>授权码</p>"
          }
        ]
      }
    },
    "success": {
      "examples": [
        {
          "title": "Response-Example",
          "content": "{\n  \"errCode\": 200,\n  \"errMsg\": \"Success OK\",\n  \"cpuUsedPercent\": \"6%\", //CPU使用率\n  \"virtualMemory\": {\n    \"total\": \"8079MB\", //总内存\n    \"available\": \"2565MB\", //当前可用内存\n    \"used\": \"5514MB\", //当前已使用内存\n    \"usedPercent\": \"68%\" //当前内存使用率\n  },\n  \"network\": {\n    \"uploadSpeed\": \"0KB/s\", //上传速度\n    \"downloadSpeed\": \"0KB/s\" //下载速度\n  },\n  \"deviceInfo\": {\n    \"totalCount\": 0, //设备总数\n    \"onlineCount\": 0 //设备在线总数\n  },\n  \"channelInfo\": {\n    \"totalCount\": 0, //通道总数\n    \"onlineCount\": 0 //通道在线总数\n  }\n}",
          "type": "json"
        }
      ]
    },
    "filename": "controllers/system.go",
    "groupTitle": "系统接口"
  }
] });
