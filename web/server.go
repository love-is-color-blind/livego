package web

// 将 RTSP 利用 FFmpeg 转换
import (
	"encoding/json"
	"github.com/gwuhaolin/livego/web/core"
	"github.com/gwuhaolin/livego/web/tools"
	"net"
	"net/http"
	"strings"
)

func Serve(l net.Listener) error {
	mux := http.NewServeMux()

	mux.HandleFunc("/list", List)
	mux.HandleFunc("/add", Add)
	mux.HandleFunc("/remove", Remove)

	http.Serve(l, mux)

	loadRtspFormDisk()

	return nil
}

var converter = core.NewRtspConverter()

func List(w http.ResponseWriter, req *http.Request) {

	list := converter.GetAll()
	resp, _ := json.Marshal(list)

	w.Write(resp)
}
func Add(w http.ResponseWriter, req *http.Request) {
	rtsp := req.FormValue("rtsp")
	var body = ""

	if len(rtsp) != 0 {
		name := converter.Add(rtsp)
		if name != "" {
			saveRtspToDisk()
			ip := req.Host
			ip = ip[0:strings.LastIndex(ip, ":")]
			var flv = "http://" + ip + ":7001/live/" + name + ".flv"
			var hls = "http://" + ip + ":7002/live/" + name + ".m3u8"
			var rtmp = "rtmp://" + ip + ":1935/live/" + name
			body = flv + "\r\n" + hls + "\r\n" + rtmp
		}
	}
	if body == "" {
		body = "false"
	}
	w.Write([]byte(body))
}
func Remove(w http.ResponseWriter, req *http.Request) {
	rtsp := req.FormValue("rtsp")
	var body = ""

	if len(rtsp) != 0 {
		name := converter.Remove(rtsp)
		if name {
			saveRtspToDisk()
			body = "true"
		}
	}
	if body == "" {
		body = "false"
	}
	w.Write([]byte(body))
}

func getFileName() string {
	if tools.IsWindows() {
		return "d:/love-db.json"
	}

	return "/home/love-db.json"
}

// 持久化
func saveRtspToDisk() {
	rtspList := converter.GetAll()
	tools.WriteToFile(getFileName(), rtspList)
}

func loadRtspFormDisk() {
	rtspList := tools.ReadFromFile(getFileName())

	rtspCount := len(rtspList)
	for i := 0; i < rtspCount; i++ {
		rtsp := rtspList[i]
		converter.Add(rtsp)
	}
}
