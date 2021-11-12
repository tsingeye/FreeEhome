### 录像查询
- 请求
```json
<?xml version="1.0" encoding="GB2312" ?>
<PPVSPMessage>
<Version>2.5</Version>
<Sequence>6194</Sequence>
<CommandType>REQUEST</CommandType>
<Method>QUERY</Method>
<Command>QUERYRECORDEDFILES</Command>
<Params>
    <Channel>2</Channel>
    <FileType>255</FileType>
    <StartTime>2021-09-15 00:00:00</StartTime>
    <StopTime>2021-09-15 23:59:59</StopTime>
    <ChunkSize>8</ChunkSize>
    <StartIdx>0</StartIdx>
</Params>

```
- 回复
```json

<?xml version="1.0" encoding="GB2312"?>
<PPVSPMessage>
<Version>2.0</Version>
<Sequence>6196</Sequence>
<CommandType>RESPONSE</CommandType>
<WhichCommand>QUERYRECORDEDFILES</WhichCommand>
<Status>200</Status>
<Description>OK</Description>
<Params>
    <RecordedFile>ch0002_00000000009000000 2021-09-15 11:08:46 2021-09-15 11:45:24 1065141708 0</RecordedFile>
    <RecordedFile>ch0002_00000000010000000 2021-09-15 11:45:24 2021-09-15 11:47:56 76423324 0</RecordedFile>
    <RecordedFile>ch0002_00000000010000100 2021-09-15 11:49:36 2021-09-15 12:03:34 399183424 0</RecordedFile>
</Params>
</PPVSPMessage>

```
### 云台控制
    start 发送三次，设备回复1次
    stop  发送三次，设备回复1次

- 请求 
```json
<?xml version="1.0" encoding="GB2312" ?>
<PPVSPMessage>
    <Version>2.5</Version>
    <Sequence>6256</Sequence>
    <CommandType>REQUEST</CommandType>
    <Method>CONTROL</Method>
    <Command>PTZCONTROL</Command>
    <Params>
        //设备的通道ID
        <Channel>2</Channel>
        //LEFT RIHGT UP DOWN
        <PTZCmd>LEFT</PTZCmd>
        //Start Stop
        <Action>Start</Action>
        //1-7
        <Speed>4</Speed>
    </Params>
</PPVSPMessage>
```
- 回复
```json
<?xml version="1.0" encoding="GB2312"?>
<PPVSPMessage>
    <Version>2.0</Version>
    <Sequence>6256</Sequence>
    <CommandType>RESPONSE</CommandType>
    <WhichCommand>PTZCONTROL</WhichCommand>
    <Status>401</Status>
    <Description>System Oper Faild.</Description>
    <Params/>
</PPVSPMessage>
```


### GPS定时上报

- 接收
```json
<?xml version="1.0" encoding="UTF-8"?>
<PPVSPMessage>
    <Version>2.0</Version>
    <Sequence>1137</Sequence>
    <CommandType>REQUEST</CommandType>
    <Command>GPS</Command>
    <Params>
        //设备ID
        <DeviceID>123457</DeviceID>
        //上报时间
        <Time>2021-08-23 16:35:35</Time>
        //方向
        <DivisionEW>E</DivisionEW>
        //经度 范围：0-（180*3600*100）
        //转换公式：实际度数*3600*100+实际分数*60*100+实际秒数*100
        <Longitude>43713595</Longitude>
        <DivisionNS>N</DivisionNS>
        //纬度
        //转换公式：实际度数*3600*100+实际分数*60*100+实际秒数*100
        <Latitude>11214118</Latitude>
        //方向 0表示正北
        <Direction>0</Direction>
        //速度 cm/h
        <Speed>0</Speed>
        //星数
        <Satellites>7</Satellites>  
        //精度
        <Precision>13</Precision>
        //高度
        <Height>2160</Height>
        <RetransFlag>1</RetransFlag>
        <NeedsResponse>1</NeedsResponse>
        <TimeZone>+08:00</TimeZone>
        <Remark>test/debug</Remark>
    </Params>
</PPVSPMessage>
```

- 回复
```json
<?xml version="1.0" encoding="GB2312" ?>
<PPVSPMessage>
    <Version>2.5</Version>
    <Sequence>1137</Sequence>
    <CommandType>RESPONSE</CommandType>
    <WhichCommand>GPS</WhichCommand>
    <Status>200</Status>
    <Description>OK</Description>
    <Params>
        <RetransFlag>1</RetransFlag>
    </Params>
</PPVSPMessage>

```