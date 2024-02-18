import { credentials } from "@grpc/grpc-js";
import { AccountsClient } from "../../../grpcGW/accounts";

export const accountsService = new AccountsClient("0.0.0.0:50051", credentials.createInsecure());