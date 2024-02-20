import { credentials } from "@grpc/grpc-js";
import { AuthClient } from "../../../grpcGW/auth";

export const authService = new AuthClient("0.0.0.0:50051", credentials.createInsecure());