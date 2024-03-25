import type { TradeClient as ITradeClient } from "./contracts/trade";
import { TradeClient } from "./contracts/trade";
import { createService } from "./utils";

export const tradeService = createService<ITradeClient>(TradeClient);
