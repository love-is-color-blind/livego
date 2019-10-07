package core

import (
	"log"
	"os/exec"
	"regexp"
	"time"
)

type RtspConverter struct {
	// map
	data map[string]*exec.Cmd
}

func NewRtspConverter() *RtspConverter {
	converter := &RtspConverter{
		data: make(map[string]*exec.Cmd),
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
	return true
}

func (c RtspConverter) getProcess(rtsp string) *exec.Cmd {
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
		r, _ := regexp.Compile("[:/@\\._]")
		name := r.ReplaceAllString(rtsp, "")

		// "ffmpeg -rtsp_transport tcp -i \"" + rtsp + "\" -vcodec copy -acodec copy -f flv  rtmp://localhost:1935/live/" + name
		log.Println(rtsp)

		cmd := exec.Command("ffmpeg",
			"-y",
			"-rtsp_transport", "tcp",
			"-i", rtsp,
			"-vcodec", "copy",
			"-acodec", "aac",

			"-f", "flv",
			"rtmp://localhost:1935/live/"+name)

		c.data[rtsp] = cmd

		cmd.Start()

		//	Tools.ProcessConsole(process).start();
		return name
	}
	return ""
}

func (c RtspConverter) stop(rtsp string) {
	process := c.getProcess(rtsp)
	if process != nil {
		process.Process.Release()
		process.Process.Kill()
	}

}

func (c RtspConverter) restart(rtsp string) {
	if rtsp == "" {
		return
	}
	log.Println(rtsp + " restart ...")
	c.stop(rtsp)
	c.start(rtsp)
}

func (c RtspConverter) KeepAlive() {
	for {
		rtspList := c.GetAll()
		for _, rtsp := range rtspList {
			process := c.getProcess(rtsp)
			if process.ProcessState != nil {
				log.Println(rtsp + " 不是活动状态，正在重启")
				//c.restart(rtsp)
			}

			time.Sleep(time.Duration(5) * time.Second)
		}
	}
}
