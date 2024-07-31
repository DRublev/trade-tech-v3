import logger from "../../logger";
import { pingService } from "../grpc/ping"

const PING_INTERVAL = 9 * 1000;
export const ping = () => {
    setInterval(() => {
        pingService.ping({}).catch(e => {
            logger.error('ping error ' + e);
        });
    }, PING_INTERVAL)
}