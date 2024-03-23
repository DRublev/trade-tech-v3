import { credentials } from "@grpc/grpc-js";
import { MarketDataClient } from "../../../contracts/marketData";

const channelOptions = {
    // Send keepalive pings every 10 seconds, default is 2 hours.
    // 'grpc.keepalive_time_ms': 10 * 1000,
    // Keepalive ping timeout after 5 seconds, default is 20 seconds.
    'grpc.keepalive_timeout_ms': 15 * 1000,
    // Allow keepalive pings when there are no gRPC calls.
    // 'grpc.keepalive_permit_without_calls': 1,
};

export const marketdataService = new MarketDataClient("0.0.0.0:50051", credentials.createInsecure(), channelOptions);