# trade-tech-v3

protoc -I protobuf protobuf/test/*.proto --go_out=./server/proto/ --go_opt=paths=source_relative --go-grpc_out=./server/proto/ --go-grpc_opt=paths=source_relative