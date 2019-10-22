package main

import (
	"flag"
	"github.com/gwuhaolin/livego/protocol/hls"
	"github.com/gwuhaolin/livego/protocol/httpflv"
	"github.com/gwuhaolin/livego/protocol/rtmp"
	"github.com/gwuhaolin/livego/web"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strings"
	"time"
)

var (
	version  = "master"
	rtmpAddr = flag.String("rtmp-addr", ":1935", "RTMP server listen address")

	hlsAddr = flag.String("hls-addr", ":7002", "HLS server listen address")
	webAddr = flag.String("web-addr", ":7777", "HTTP Web interface server listen address")
)

func init() {
	log.SetFlags(log.Lshortfile | log.Ltime | log.Ldate)
	flag.Parse()
}

func startRtmp(stream *rtmp.RtmpStream, hlsServer *hls.Server) {
	rtmpListen, err := net.Listen("tcp", *rtmpAddr)
	if err != nil {
		log.Fatal(err)
	}

	var rtmpServer *rtmp.Server

	if hlsServer == nil {
		rtmpServer = rtmp.NewRtmpServer(stream, nil)
		log.Printf("hls server disable....")
	} else {
		rtmpServer = rtmp.NewRtmpServer(stream, hlsServer)
		log.Printf("hls server enable....")
	}

	defer func() {
		if r := recover(); r != nil {
			log.Println("RTMP server panic: ", r)
		}
	}()
	log.Println("RTMP Listen On", *rtmpAddr)
	rtmpServer.Serve(rtmpListen)
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			log.Println("livego panic: ", r)
			time.Sleep(1 * time.Second)
		}
	}()
	log.Println("start livego, version", version)

	stream := rtmp.NewRtmpStream()

	opListen, err := net.Listen("tcp", *webAddr)
	if err != nil {
		log.Fatal(err)
		return
	}
	mux := http.NewServeMux()

	opServer := web.NewServer(stream, *rtmpAddr)
	opServer.AddStatApiUrl(mux)
	web.AddRTSPApiUrl(mux)

	flvServer := httpflv.NewServer(stream)
	hlsServer := hls.NewServer()
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		path := request.URL.Path
		isFlv := strings.HasSuffix(path, ".flv")
		isHls := strings.HasSuffix(path, ".crossdomain.xml") || strings.HasSuffix(path, ".m3u8") || strings.HasSuffix(path, ".ts")
		if isFlv {
			flvServer.HandleConn(writer, request)
		} else if isHls {
			hlsServer.Handle(writer, request)
		} else {
			bytes, _ := ioutil.ReadFile("index.html")
			writer.Write(bytes)
		}
	})

	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Println("HTTP-API server panic: ", r)
			}
		}()
		log.Println("HTTP-API listen On", *webAddr)

		http.Serve(opListen, mux)
	}()

	startRtmp(stream, hlsServer)
}
