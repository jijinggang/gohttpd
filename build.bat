SET GOARCH=amd64

SET GOOS=windows
go build -ldflags "-s -w"

rem build linux64
SET GOOS=linux
go build -ldflags "-s -w"

