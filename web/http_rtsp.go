package web

// 将 RTSP 利用 FFmpeg 转换
import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
)

func AddRTSPApiUrl(mux *http.ServeMux) {
	loadRtspFormDisk()

	mux.HandleFunc("/list", list)
	mux.HandleFunc("/add", add)
	mux.HandleFunc("/remove", remove)
}

var converter = NewRtspConverter()

func list(w http.ResponseWriter, req *http.Request) {
	list := converter.GetAll()
	ip := getIp(req)
	port := getPort(req)
	var rsList []StreamInfo

	for _, rtsp := range list {
		info := GetInfoByRTSP(rtsp, ip, port)
		rsList = append(rsList, info)
	}

	bytes, _ := json.Marshal(rsList)
	w.Write(bytes)

}
func add(w http.ResponseWriter, req *http.Request) {
	rtsp := req.FormValue("rtsp")
	if len(rtsp) != 0 {
		name := converter.Add(rtsp)
		if name != "" {
			saveRtspToDisk()
			ip := getIp(req)
			port := getPort(req)
			addressList := GetInfoByRTSP(rtsp, ip, port)
			bytes, e := json.Marshal(addressList)
			if e != nil {
				w.Write([]byte("false"))
			} else {
				w.Write(bytes)
			}
		}
	}
}

func getIp(req *http.Request) string {
	ip := req.Host
	ip = ip[0:strings.LastIndex(ip, ":")]
	return ip
}
func getPort(req *http.Request) string {
	ip := req.Host
	ip = ip[strings.LastIndex(ip, ":")+1:]
	return ip
}
func remove(w http.ResponseWriter, req *http.Request) {
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
