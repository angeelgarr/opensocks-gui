#!bin/bash
export GO111MODULE=on
echo "start building..."

GOOS=linux GOARCH=amd64 go build -o ./bin/opensocks-gui ./main.go
#GOOS=linux GOARCH=arm64 go build -o ./bin/opensocks-gui-linux-arm64 ./main.go
#GOOS=darwin GOARCH=amd64 go build -o ./bin/opensocks-gui-macos ./main.go
#GOOS=windows GOARCH=amd64 go build -o ./bin/opensocks-gui.exe ./main.go

echo "done!"
