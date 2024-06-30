
const grpc = require("@grpc/grpc-js");
var protoLoader = require("@grpc/proto-loader");

const options = {
    keepCase: true,
    longs: String,
    enums: String,
    defaults: true,
    oneofs: true,
};

const AUTH_PROTO_PATH = "../protobuf/auth.proto";
var authPackageDefinition = protoLoader.loadSync(AUTH_PROTO_PATH, options);
const authProto = grpc.loadPackageDefinition(authPackageDefinition);

const ACCOUNTS_PROTO_PATH = "../protobuf/accounts.proto";
var accountsPackageDefinition = protoLoader.loadSync(ACCOUNTS_PROTO_PATH, options);
const accountsProto = grpc.loadPackageDefinition(accountsPackageDefinition);

const server = new grpc.Server();

server.addService(authProto.auth.Auth.service, {
    HasToken: (_, callback) => {
        callback(null, { HasToken: true });
    },
});
server.addService(accountsProto.accounts.Accounts.service, {
    GetAccount: (_, callback) => {
        callback(null, { AccountId: "mocked-account-id" });
    },
});

server.bindAsync(
    "127.0.0.1:50051",
    grpc.ServerCredentials.createInsecure(),
    (error, port) => {
        console.log("Server running at http://127.0.0.1:50051");
        server.start();
    }
);