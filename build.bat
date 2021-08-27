REM build windows
go build -ldflags "-s -w"


REM build linux
SET GOOS=linux
SET GOARCH=amd64
go build -ldflags "-s -w"

