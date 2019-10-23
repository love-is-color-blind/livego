# 将RTSP 协议转换为如下
#### 支持的传输协议
- [x] RTMP
- [x] AMF
- [x] HLS
- [x] HTTP-FLV

#### 支持的容器格式
- [x] FLV
- [x] TS

#### 支持的编码格式
- [x] H264
- [x] AAC
- [x] MP3



## 使用
docker hub 中搜索 rtsp-live-stream
可以持久化 /db.txt 防止重启后数据丢失
### 转推模式，可以从一台服务器推到另一台服务器  
    设置环境变量，LIVE_STREAM_PUSH_TARGET_SERVER=另一台服务器IP
    例如内网RTSP接口，推送到公网
    




## 测试
# 获得测试的RTSP流 ？
使用vlc，可以将mp4，或者其他输入源串流为RTSP

