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
	"time"
)

var (
	version  = "master"
	rtmpAddr = flag.String("rtmp-addr", ":1935", "RTMP server listen address")

	hlsAddr = flag.String("hls-addr", ":7002", "HLS server listen address")

	httpFlvAddr = flag.String("httpflv-addr", ":7001", "HTTP-FLV server listen address")
	webAddr     = flag.String("web-addr", ":7777", "HTTP Web interface server listen address")
)

func init() {
	log.SetFlags(log.Lshortfile | log.Ltime | log.Ldate)
	flag.Parse()
}

func startHls() *hls.Server {
	hlsListen, err := net.Listen("tcp", *hlsAddr)
	if err != nil {
		log.Fatal(err)
	}

	hlsServer := hls.NewServer()
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Println("HLS server panic: ", r)
			}
		}()
		log.Println("HLS listen On", *hlsAddr)
		hlsServer.Serve(hlsListen)
	}()
	return hlsServer
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

func startHTTPFlv(stream *rtmp.RtmpStream) {
	flvListen, err := net.Listen("tcp", *httpFlvAddr)
	if err != nil {
		log.Fatal(err)
	}

	hdlServer := httpflv.NewServer(stream)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Println("HTTP-FLV server panic: ", r)
			}
		}()
		log.Println("HTTP-FLV listen On", *httpFlvAddr)
		hdlServer.Serve(flvListen)
	}()
}

func startWeb(stream *rtmp.RtmpStream) {
	if *webAddr != "" {
		opListen, err := net.Listen("tcp", *webAddr)
		if err != nil {
			log.Fatal(err)
		}
		mux := http.NewServeMux()

		opServer := web.NewServer(stream, *rtmpAddr)
		opServer.AddOperaUrl(mux)

		web.AddRTSPUrl(mux)

		mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
			bytes, _ := ioutil.ReadFile("index.html")
			writer.Write(bytes)
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
	}
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
	hlsServer := startHls()
	startHTTPFlv(stream)
	startWeb(stream)
	startRtmp(stream, hlsServer)

}
