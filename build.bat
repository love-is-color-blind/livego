:: windows 打包
:: go build rtsp-live-stream.go

:: linux 打包
SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=amd64
go build rtsp-live-stream.go
pause