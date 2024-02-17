# trade-tech-v3

GO
protoc -I protobuf protobuf/*.proto --go_out=./server/grpcGW/ --go_opt=paths=import --go-grpc_out=./server/grpcGW/ --go-grpc_opt=paths=import


JS
protoc --plugin=protoc-gen-ts_proto=".\\client\\node_modules\\.bin\\protoc-gen-ts_proto.cmd" --ts_proto_out=./client/grpcGW --ts_proto_opt=outputServices=grpc-js --ts_proto_opt=esModuleInterop=true -I ./protobuf ./protobuf/*.proto