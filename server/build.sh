#!/bin/bash
cd ../server

GOOS=windows GOARCH=amd64 go build -o ./resources/app/go-binaries/app-binary-windows.exe -ldflags "-X main.envFromBuild=PROD -X main.secretFromBuild=trade-tech-secret-for-encryption"
GOOS=darwin GOARCH=amd64 go build -o ./resources/app/go-binaries/app-binary-macos -ldflags "-X main.envFromBuild=PROD -X main.secretFromBuild=trade-tech-secret-for-encryption"
GOOS=linux GOARCH=amd64  go build -o ./resources/app/go-binaries/app-binary-linux -ldflags "-X main.envFromBuild=PROD -X main.secretFromBuild=trade-tech-secret-for-encryption"