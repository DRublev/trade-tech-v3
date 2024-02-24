import { credentials } from "@grpc/grpc-js";
import { SharesClient, SharesService } from "../../../grpcGW/instruments";

export const instrumentsService = new SharesClient("0.0.0.0:50051", credentials.createInsecure());