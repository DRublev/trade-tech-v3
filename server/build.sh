#!/bin/bash
cd ../server

GOOS=windows GOARCH=amd64 go build -o ./resources/app/go-binaries/app-binary-windows.exe
GOOS=darwin GOARCH=amd64 go build -o ./resources/app/go-binaries/app-binary-macos
GOOS=linux GOARCH=amd64  go build -o ./resources/app/go-binaries/app-binary-linux