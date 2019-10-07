package web

// 将 RTSP 利用 FFmpeg 转换
import (
	"encoding/json"
	"github.com/gwuhaolin/livego/web/core"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strings"
)

func Serve(l net.Listener) error {
	loadRtspFormDisk()

	mux := http.NewServeMux()
	mux.HandleFunc("/list", List)
	mux.HandleFunc("/add", Add)
	mux.HandleFunc("/remove", Remove)

	error := http.Serve(l, mux)
	if error != nil {
		log.Println(error)
	}
	return error
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

// 持久化
func saveRtspToDisk() {
	rtspList := converter.GetAll()
	writeToFile("db.txt", rtspList)
}

func loadRtspFormDisk() {
	rtspList := readFromFile("db.txt")

	rtspCount := len(rtspList)
	for i := 0; i < rtspCount; i++ {
		rtsp := rtspList[i]
		converter.Add(rtsp)
	}
}

func writeToFile(fileName string, lines []string) {
	data := []byte(strings.Join(lines, "\n"))
	ioutil.WriteFile(fileName, data, 777)
}

func readFromFile(fileName string) []string {
	data, error := ioutil.ReadFile(fileName)
	if error != nil {
		return nil
	}
	return strings.Split(string(data), "\n")
}
