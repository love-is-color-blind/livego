package web

// 将rtsp 转换成rtmp
import (
	"crypto/md5"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

type RtspConverter struct {
	// map
	data map[string]*CmdInfo
}

type CmdInfo struct {
	cmd    *exec.Cmd
	rtsp   string
	writer MyWriter
}

type StreamInfo struct {
	Key     string `json:"key"`
	Rtsp    string `json:"rtsp"`
	Rtmp    string `json:"rtmp"`
	HttpFlv string `json:"httpFlv"`
	Hls     string `json:"hls"`
}

func GetInfo(rtsp string, ip string, port string) StreamInfo {
	name := getRTSPKey(rtsp)

	host := getTargetHost()
	if host != "localhost" {
		ip = host
	}

	var flv = "http://" + ip + ":" + port + "/live/" + name + ".flv"
	var hls = "http://" + ip + ":" + port + "/live/" + name + ".m3u8"
	var rtmp = "rtmp://" + ip + ":1935/live/" + name

	return StreamInfo{
		Key:     name,
		Rtsp:    rtsp,
		Rtmp:    rtmp,
		HttpFlv: flv,
		Hls:     hls,
	}
}

func NewRtspConverter() *RtspConverter {
	converter := &RtspConverter{
		data: make(map[string]*CmdInfo),
	}
	go func() {
		converter.KeepAlive()
	}()

	return converter
}

func (c *RtspConverter) Add(rtsp string) string {
	if rtsp != "" {
		c.stop(rtsp)
		return c.start(rtsp)
	}
	return ""
}

func (c RtspConverter) Remove(rtsp string) bool {
	c.stop(rtsp)
	delete(c.data, rtsp)
	return true
}

func (c RtspConverter) getProcess(rtsp string) *CmdInfo {
	cmd := c.data[rtsp]
	return cmd
}

func (c RtspConverter) GetAll() []string {
	var rtspList []string
	for k := range c.data {
		rtspList = append(rtspList, k)
	}
	return rtspList
}

// 开启一个rtsp 转换服务， 返回 flv hls rtmp 访问地址

func (c RtspConverter) start(rtsp string) string {
	if rtsp != "" {
		name := getRTSPKey(rtsp)
		// "ffmpeg -rtsp_transport tcp -i \"" + rtsp + "\" -vcodec copy -acodec aac -f flv  rtmp://localhost:1935/live/" + name
		log.Println(rtsp)

		cmd := exec.Command("ffmpeg",
			"-y",
			"-rtsp_transport", "tcp",
			"-i", rtsp,
			"-vcodec", "copy",
			"-acodec", "aac",
			"-ar", "44100",
			"-f", "flv",
			"rtmp://"+getTargetHost()+":1935/live/"+name)
		log.Println(strings.Join(cmd.Args, " "))

		writer := MyWriter{converter: c, rtsp: rtsp, lastWriteTime: time.Now().Unix()}
		cmd.Stderr = writer
		cmd.Start()

		c.data[rtsp] = &CmdInfo{
			cmd:    cmd,
			rtsp:   rtsp,
			writer: writer,
		}
		return name
	}
	return ""
}

func getRTSPKey(rtsp string) string {
	data := []byte(rtsp)
	has := md5.Sum(data)
	md5str1 := fmt.Sprintf("%x", has) //将[]byte转成16进制
	fmt.Println(md5str1)
	return md5str1
}

func (c RtspConverter) stop(rtsp string) {
	process := c.getProcess(rtsp)
	if process != nil {
		process.cmd.Process.Release()
		process.cmd.Process.Kill()
	}
}

func (c RtspConverter) restart(rtsp string) {
	if rtsp == "" {
		return
	}
	log.Println("############## restart ", rtsp)
	c.stop(rtsp)
	c.start(rtsp)
}

func (c RtspConverter) KeepAlive() {
	for {
		rtspList := c.GetAll()
		for _, rtsp := range rtspList {
			process := c.getProcess(rtsp)

			timeout := time.Now().Unix() - process.writer.lastWriteTime
			if process.cmd.ProcessState != nil || timeout > 60 {
				log.Println("################ Not Alive , Restarting #######################")
				c.restart(rtsp)
			}

			time.Sleep(5 * time.Second)
		}
	}
}

type MyWriter struct {
	converter     RtspConverter
	rtsp          string
	lastWriteTime int64 // seconds
}

func (w MyWriter) Write(p []byte) (n int, err error) {
	w.lastWriteTime = time.Now().Unix()
	str := string(p)
	log.Println(str)
	return len(p), nil
}

func getTargetHost() string {
	host := os.Getenv("LIVE_STREAM_PUSH_TARGET_SERVER")
	if host == "" {
		host = "localhost"
	}
	return host
}
