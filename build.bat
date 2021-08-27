go build -ldflags "-s -w"

rem build linux64
SET GOOS=linux
SET GOARCH=amd64
go build -ldflags "-s -w"

