import { credentials } from "@grpc/grpc-js";
import { InstrumentsClient } from "../../../grpcGW/instruments";

export const accountsService = new InstrumentsClient("0.0.0.0:50051", credentials.createInsecure());