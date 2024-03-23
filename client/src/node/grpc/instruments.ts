import { credentials } from "@grpc/grpc-js";
import { SharesClient, SharesService } from "../../../contracts/shares";

export const sharesService = new SharesClient("0.0.0.0:50051", credentials.createInsecure());