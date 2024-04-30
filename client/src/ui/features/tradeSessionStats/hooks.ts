import { useCallback, useEffect, useRef, useState } from "react";
import { OrderOperations, OrderState } from "../../../types";
import { useCurrentInstrument } from "../../utils/useCurrentInstrumentId";
import { useOrders } from "../space/hooks";
import { useLogger } from "../../hooks";

const calcedTrades: Record<string, boolean> = {}
export const useTradeSessionStats = () => {
    const [instrument] = useCurrentInstrument();
    const [turnover, setTurnover] = useState(0);
    const [profit, setProfit] = useState(0);
    const [tradesAmount, setTradesAmount] = useState(0);
    const buyPrices = useRef([]);
    const logger = useLogger({ component: 'useTradeSessionStats' });

    const handleOrderStateChange = useCallback((orderState: OrderState) => {
        logger.trace('Got info about new order', { isPartiallyExecuted: orderState.lotsExecuted != orderState.lotsRequested });
        // Частично исполненная зхаявка, пока хз как их считать
        // if (orderState.lotsExecuted != orderState.lotsRequested) return;

        if (calcedTrades[orderState.id]) return;
        calcedTrades[orderState.id] = true

        const price = orderState.price;

        let profit = 0;
        if (orderState.operationType == OrderOperations.Buy) {
            for (let i = 0; i < orderState.lotsExecuted; i++) {
                buyPrices.current.push(price);
            }
        } else {
            // Продается всегда акция, купленная раньше всех, поэтому здесь FIFO
            for (let i = 0; i < orderState.lotsExecuted; i++) {
                const headBuyPrice = buyPrices.current.shift();
                buyPrices.current.push(price);
                if (!headBuyPrice) {
                    continue;
                }
                profit += price - headBuyPrice;
            }
        }

        setTurnover(t => t + price);
        setProfit(p => p + profit);
        setTradesAmount(a => a + 1);
    }, [buyPrices.current]);

    useOrders(handleOrderStateChange, instrument);

    return [turnover, profit, tradesAmount];
};

type OrderLog = OrderState;
const calcedLogs: Record<string, boolean> = {};
export const useTradeLogs = () => {
    const [instrument] = useCurrentInstrument();
    const [logs, setLogs] = useState<OrderLog[]>([]);

    const handleOrderStateChange = (orderState: OrderState) => {
        if (calcedLogs[orderState.id]) return;
        calcedLogs[orderState.id] = true;
        
        setLogs([...logs, orderState]);
    };
    useOrders(handleOrderStateChange, instrument);


    useEffect(() => {
        setLogs([]);
    }, [instrument]);

    return logs;
};

