import type { FC } from 'react';
import React from 'react';
import { useTradeSessionStats } from './hooks';
import { Card, Heading, Text, Flex } from '@radix-ui/themes';
import { DoubleArrowUpIcon, ListBulletIcon, ReloadIcon } from '@radix-ui/react-icons';

type Props = never;

export const TradeStats: FC<Props> = () => {
    const [turnover, profit, tradesAmount] = useTradeSessionStats();

    return (<Card>
        <Heading mb="1">Статистика</Heading>
        <Flex align="center" mb="1" gap="2">
            <ReloadIcon />
            <Text>Оборот</Text>
            {turnover.toFixed(2)}
        </Flex>
        <Flex align="center" mb="1" gap="2">
            <DoubleArrowUpIcon />
            <Text>Профит</Text>
            <Text color={profit >= 0 ? 'green' : 'red'}>{profit.toFixed(2)}</Text>
        </Flex>
        <Flex align="center" mb="1" gap="2">
            <ListBulletIcon />
            <Text>Количество сделок</Text>
            <Text>{tradesAmount}</Text>

        </Flex>
    </Card>);
};
