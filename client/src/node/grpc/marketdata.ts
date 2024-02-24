import { credentials } from "@grpc/grpc-js";
import { MarketDataClient } from "../../../grpcGW/marketData";

export const marketdataService = new MarketDataClient("0.0.0.0:50051", credentials.createInsecure());