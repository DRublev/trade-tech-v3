# trade-tech-v3

GO
protoc -I protobuf protobuf/test/*.proto --go_out=./server/proto/ --go_opt=paths=source_relative --go-grpc_out=./server/proto/ --go-grpc_opt=paths=source_relative

JS
protoc --plugin=protoc-gen-ts_proto=".\\node_modules\\.bin\\protoc-gen-ts_proto.cmd" --ts_proto_out=./protobuf --ts_proto_opt=outputServices=grpc-js --ts_proto_opt=esModuleInterop=true -I ../protobuf ../protobuf/*.proto