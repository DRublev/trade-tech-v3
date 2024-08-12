#!/bin/bash
cd ../server

mkdir -p ./resources/app/go-binaries
chown -R $(whoami) ./resources/app/go-binaries

GOOS=windows GOARCH=amd64 CGO_ENABLED=1 GOGCCFLAGS='' CC=gcc go build -o ./resources/app/go-binaries/app-binary-windows.dll -buildmode=c-shared -ldflags "-X main.envFromBuild=PROD -X main.secretFromBuild=trade-tech-secret-for-encryption"
GOOS=darwin GOARCH=arm64 CGO_ENABLED=1 go build -o ./resources/app/go-binaries/app-binary-macos.so -buildmode=c-shared -ldflags "-X main.envFromBuild=PROD -X main.secretFromBuild=trade-tech-secret-for-encryption"
# GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build -o ./resources/app/go-binaries/app-binary-linux.so -buildmode=c-shared -ldflags "-X main.envFromBuild=PROD -X main.secretFromBuild=trade-tech-secret-for-encryption"

# GOOS=windows GOARCH=amd64 go build -o ./resources/app/go-binaries/app-binary-windows.exe -ldflags "-X main.envFromBuild=PROD -X main.secretFromBuild=trade-tech-secret-for-encryption"
# GOOS=darwin GOARCH=amd64 go build -o ./resources/app/go-binaries/app-binary-macos -ldflags "-X main.envFromBuild=PROD -X main.secretFromBuild=trade-tech-secret-for-encryption"
# GOOS=linux GOARCH=amd64  go build -o ./resources/app/go-binaries/app-binary-linux -ldflags "-X main.envFromBuild=PROD -X main.secretFromBuild=trade-tech-secret-for-encryption"