package core

import (
	"log"
	"os/exec"
	"regexp"
)

type RtspConverter struct {
	// map
	data map[string]*exec.Cmd
}

func NewRtspConverter() *RtspConverter {
	return &RtspConverter{}
}

func (c *RtspConverter) Add(rtsp string) string {
	c.stop(rtsp)
	return c.start(rtsp)
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

	slice1 := make([]string, 5)
	for k := range c.data {
		slice1 = append(slice1, k)
	}
	return slice1
}

// 开启一个rtsp 转换服务， 返回 flv hls rtmp 访问地址

func (c RtspConverter) start(rtsp string) string {
	if rtsp != "" {
		r, _ := regexp.Compile("[:/@\\._]")
		name := r.ReplaceAllString(rtsp, "")

		fullCommand := "-rtsp_transport tcp -i \"" + rtsp + "\" -vcodec copy -acodec copy -f flv  rtmp://localhost:1935/live/" + name
		log.Println(rtsp)
		log.Println(fullCommand)

		command := exec.Command("ffmpeg",
			"-y",
			"-rtsp_transport", "tcp",
			"-i", "\""+rtsp+"\"",
			"-vcodec", "copy",
			"-acodec", "copy",
			"-f", "flv",
			"rtmp://localhost:1935/live/"+name)
		c.data[rtsp] = command

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
