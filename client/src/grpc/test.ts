import { credentials } from "@grpc/grpc-js";
import { TestClient } from "../../protobuf/test";

export const testService = new TestClient("0.0.0.0:50051", credentials.createInsecure());
