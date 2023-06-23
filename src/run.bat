set GOOS=linux
set GOARCH=amd64
del buildevader
go build -ldflags="-s -w" -buildmode=exe buildevader.go
del d:\TestCases\buildevader
copy buildevader d:\TestCases