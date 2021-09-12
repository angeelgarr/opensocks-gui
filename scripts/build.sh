#!bin/bash
export GO111MODULE=on

GOOS=linux GOARCH=amd64 go build -o ./bin/opensocks-gui-linux-amd64 ./main.go
GOOS=darwin GOARCH=amd64 CGO_ENABLED=1 CC=o64-clang CXX=o64-clang++ go build -ldflags -s -o ./bin/opensocks-gui-darwin-amd64 ./main.go
GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc CXX=x86_64-w64-mingw32-g++ go build -ldflags "-H windowsgui"  -o ./bin/opensocks-gui-amd64.exe ./main.go

echo "DONE!!!"
