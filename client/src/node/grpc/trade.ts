import { credentials } from "@grpc/grpc-js";
import { TradeClient } from "../../../grpcGW/trade";

export const tradeService = new TradeClient("0.0.0.0:50051", credentials.createInsecure());