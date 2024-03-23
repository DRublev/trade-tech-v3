# trade-tech-v3

## Подготовка к запуску

Скачать [Go](https://go.dev/doc/install)

Скачать [Node.js](https://nodejs.org/en/download)

Выпустить [токен для Тинькофф Инвестиций](https://tinkoff.github.io/investAPI/token/) (для тестирования без торговли достаточно с доступом только для чтения)
Рекомендуется завести отдельный брокерский счет

## Запуск

Склонировать проект
Для запуска клиента, нужно запустить также и сервер

### Запуск сервера 
Перейти в директорию `server`
Добавить файл `.env` и вставить в него строку `SECRET="trade-tech-secret-for-encryption"`
Запустить команду `go run .`
Запустится gRpc сервер на 50051 порту 


### Запуск клиента
Запустить сервер
Перейти в диеркторию `client`
Установить зависимости командой `npm i`
Запустить приложение `npm run start`

## Генерация protoc

Скачать [protoc](https://grpc.io/docs/protoc-installation/)

TODO: Докинуть установку плагинов для go и ts (прям командами)

Установить плагины для ts
```sh
npm i -g ts-proto
npm i -g protoc-gen-ts
```

Установить плагины для go
```sh
brew install protoc-gen-go
brew install protoc-gen-go-grpc
```

Запустить команды
GO
```sh
([[ ! -d ./server/server/contracts ]] && mkdir ./server/server/contracts) || protoc -I protobuf protobuf/*.proto --go_out=./server/server/contracts/ --go_opt=paths=import --go-grpc_out=./server/server/contracts/ --go-grpc_opt=paths=import
```

TS for Windows
```sh
protoc --plugin=protoc-gen-ts_proto=".\\client\\node_modules\\.bin\\protoc-gen-ts_proto.cmd" --ts_proto_out=./client/contracts --ts_proto_opt=outputServices=grpc-js --ts_proto_opt=esModuleInterop=true -I ./protobuf ./protobuf/*.proto
```
TS Unix
```sh
([[ ! -d ./client/contracts ]] && mkdir ./client/contracts) || protoc --ts_proto_out=./client/contracts --ts_proto_opt=outputServices=grpc-js --ts_proto_opt=esModuleInterop=true -I ./protobuf ./protobuf/*.proto
```
