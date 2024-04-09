import { useCallback, useEffect, useRef, useState } from "react";
import { OrderOperations, OrderState } from "../../../types";
import { useCurrentInstrument } from "../../utils/useCurrentInstrumentId";
import { useOrders } from "../space/hooks";


export const useTradeSessionStats = () => {
    const [instrument] = useCurrentInstrument();
    const [turnover, setTurnover] = useState(0);
    const [profit, setProfit] = useState(0);
    const [tradesAmount, setTradesAmount] = useState(0);
    const buyPrices = useRef([])

    const handleOrderStateChange = useCallback((orderState: OrderState) => {
        // Частично исполненная зхаявка, пока хз как их считать
        if (orderState.lotsExecuted != orderState.lotsRequested) return;

        const price = orderState.price * orderState.lotsExecuted;

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
export const useTradeLogs = () => {
    const [instrument] = useCurrentInstrument();
    const [logs, setLogs] = useState<OrderLog[]>([]);

    const handleOrderStateChange = (orderState: OrderState) => {
        if (orderState.lotsExecuted === orderState.lotsRequested) {
            const newLogs = logs;
            newLogs.push(orderState);
            setLogs(newLogs);
        }
    };
    useOrders(handleOrderStateChange, instrument);


    useEffect(() => {
        setLogs([]);
    }, [instrument]);

    return logs;
};

