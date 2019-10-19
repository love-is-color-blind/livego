go build rtsp-live-stream.go
SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=amd64
go build rtsp-live-stream.go
pause