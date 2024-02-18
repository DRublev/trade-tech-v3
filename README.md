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
Запустить команды
GO
```sh
protoc -I protobuf protobuf/*.proto --go_out=./server/grpcGW/ --go_opt=paths=import --go-grpc_out=./server/grpcGW/ --go-grpc_opt=paths=import
```

TS
```sh
protoc --plugin=protoc-gen-ts_proto=".\\client\\node_modules\\.bin\\protoc-gen-ts_proto.cmd" --ts_proto_out=./client/grpcGW --ts_proto_opt=outputServices=grpc-js --ts_proto_opt=esModuleInterop=true -I ./protobuf ./protobuf/*.proto
```
