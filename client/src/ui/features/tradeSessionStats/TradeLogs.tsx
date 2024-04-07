import type { FC } from "react";
import React from "react";
import { useTradeLogs } from "./hooks";
import { ArrowUpIcon, ArrowDownIcon } from "@radix-ui/react-icons";
import { Card, Flex, Heading, ScrollArea, Text, Tooltip } from "@radix-ui/themes";
import { OrderOperations } from "../../../types";

type Props = object;

const BuyLabel = () => (
    <Tooltip content="Покупка">
        <ArrowDownIcon color="#b5d2c1" />
    </Tooltip>
);
const SellLabel = () => (
    <Tooltip content="Продажа">
        <ArrowUpIcon color="#ff7f7f" />
    </Tooltip>
);

export const TradeLogs: FC<Props> = () => {
    const logs = useTradeLogs();

    return (
        <Card>
            <Heading mb="1">Сделки</Heading>
            <ScrollArea scrollbars="vertical" style={{ maxHeight: '15vh' }}>
                {logs.length === 0 && <Text>Пока нет совершенных сделок</Text>}
                {logs.map((l) => (
                    <Flex gap="5" key={l.time} align="center" p="1">
                        {l.operationType == OrderOperations.Buy ? <BuyLabel /> : <SellLabel />}
                        <Text>{l.lotsExecuted} x {l.price}₽</Text>
                        <Text>{l.price * l.lotsExecuted}₽</Text>
                    </Flex>
                ))}
            </ScrollArea>
        </Card>
    );
};
