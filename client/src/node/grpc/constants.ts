import { credentials } from "@grpc/grpc-js";

export const DEFAULT_ADDRESS = "0.0.0.0:50051";
export const DEFAULT_CREDS = credentials.createInsecure();