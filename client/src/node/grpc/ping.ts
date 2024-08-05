import { PingClient } from "./contracts/ping";
import type { PingClient as IPingClient } from "./contracts/ping";
import { createService } from "./utils";

export const pingService = createService<IPingClient>(PingClient);
