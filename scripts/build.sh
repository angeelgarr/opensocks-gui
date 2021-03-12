#!bin/bash
export GO111MODULE=on
echo "start building..."

GOOS=linux GOARCH=amd64 go build -o ./bin/opensocks-gui ./main.go
#GOOS=linux GOARCH=arm64 go build -o ./bin/opensocks-gui-linux-arm64 ./main.go
#GOOS=darwin GOARCH=amd64 go build -o ./bin/opensocks-gui-macos ./main.go
GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc CXX=x86_64-w64-mingw32-g++ go build -ldflags "-H windowsgui"  -o ./bin/opensocks-gui.exe ./main.go

echo "done!"
