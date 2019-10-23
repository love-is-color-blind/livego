package web

import (
	"encoding/json"
	"github.com/gwuhaolin/livego/av"
	"github.com/gwuhaolin/livego/protocol/rtmp"
	"github.com/gwuhaolin/livego/protocol/rtmp/rtmprelay"
	"net/http"
)

type Server struct {
	handler  av.Handler
	session  map[string]*rtmprelay.RtmpRelay
	rtmpAddr string
}

func NewServer(h av.Handler, rtmpAddr string) *Server {
	return &Server{
		handler:  h,
		session:  make(map[string]*rtmprelay.RtmpRelay),
		rtmpAddr: rtmpAddr,
	}
}

func (s *Server) AddStatApiUrl(mux *http.ServeMux) {
	mux.HandleFunc("/rtmp/statics", func(w http.ResponseWriter, r *http.Request) {
		s.GetLiveStatics(w, r)
	})
	mux.HandleFunc("/rtmp/list", func(w http.ResponseWriter, r *http.Request) {
		s.GetLiveList(w, r)
	})
}

type stream struct {
	Key             string `json:"Key"`
	Url             string `json:"Url"`
	StreamId        uint32 `json:"StreamId"`
	VideoTotalBytes uint64 `json:123456`
	VideoSpeed      uint64 `json:123456`
	AudioTotalBytes uint64 `json:123456`
	AudioSpeed      uint64 `json:123456`
}

type streams struct {
	Publishers []stream `json:"publishers"`
	Players    []stream `json:"players"`
}

func (server *Server) GetLiveStatics(w http.ResponseWriter, req *http.Request) {
	msgs, done := getStreams(server)
	if done {
		return
	}
	resp, _ := json.Marshal(msgs)
	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
}

func (server *Server) GetLiveList(w http.ResponseWriter, req *http.Request) {
	msgs, done := getStreams(server)
	if done {
		return
	}

	port := getPort(req)
	ip := getIp(req)
	var list []StreamInfo
	for _, pub := range msgs.Publishers {
		list = append(list, GetInfoByKey(ip, port, pub.Key))
	}

	resp, _ := json.Marshal(list)
	w.Write(resp)
}

func getStreams(server *Server) (*streams, bool) {
	rtmpStream := server.handler.(*rtmp.RtmpStream)
	if rtmpStream == nil {
		return nil, true
	}
	msgs := new(streams)
	for item := range rtmpStream.GetStreams().IterBuffered() {
		if s, ok := item.Val.(*rtmp.Stream); ok {
			if s.GetReader() != nil {
				switch s.GetReader().(type) {
				case *rtmp.VirReader:
					v := s.GetReader().(*rtmp.VirReader)
					msg := stream{item.Key, v.Info().URL, v.ReadBWInfo.StreamId, v.ReadBWInfo.VideoDatainBytes, v.ReadBWInfo.VideoSpeedInBytesperMS,
						v.ReadBWInfo.AudioDatainBytes, v.ReadBWInfo.AudioSpeedInBytesperMS}
					msgs.Publishers = append(msgs.Publishers, msg)
				}
			}
		}
	}
	for item := range rtmpStream.GetStreams().IterBuffered() {
		ws := item.Val.(*rtmp.Stream).GetWs()
		for s := range ws.IterBuffered() {
			if pw, ok := s.Val.(*rtmp.PackWriterCloser); ok {
				if pw.GetWriter() != nil {
					switch pw.GetWriter().(type) {
					case *rtmp.VirWriter:
						v := pw.GetWriter().(*rtmp.VirWriter)
						msg := stream{item.Key, v.Info().URL, v.WriteBWInfo.StreamId, v.WriteBWInfo.VideoDatainBytes, v.WriteBWInfo.VideoSpeedInBytesperMS,
							v.WriteBWInfo.AudioDatainBytes, v.WriteBWInfo.AudioSpeedInBytesperMS}
						msgs.Players = append(msgs.Players, msg)
					}
				}
			}
		}
	}
	return msgs, false
}
