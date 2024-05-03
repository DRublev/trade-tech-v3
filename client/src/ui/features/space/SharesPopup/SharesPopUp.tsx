import * as ToggleGroup from '@radix-ui/react-toggle-group';
import { MixerHorizontalIcon } from "@radix-ui/react-icons";
import { Box, Button, Card, Container, Flex, Text } from "@radix-ui/themes";
import React, { useCallback, useEffect, useMemo, useState } from "react";
import { PopoverWindow } from "../../../components/PopoverWindow";
import style from '../../../basicStyles.css';
import { Quatation, Share, TradingSchedule } from "../../../../node/grpc/contracts/shares";
import { useTodaysSchedules, useSharesFromStore } from '../hooks';
import { SearchInput } from '../../../components/SearchInput';
import s from './styles.css';
import { useCurrentInstrument } from '../../..//utils/useCurrentInstrumentId';
import { useLogger } from '../../../hooks';

const nanoPrecision = 1_000_000_000;
const quantToNumber = (q: Quatation | undefined): number => {
    return q ? Number(q.units + (q.nano / nanoPrecision)) : 0;
}

const isContainsWithIgnoreCase = (value: string, term: string): boolean => {
    return value.toLocaleLowerCase().includes(term.toLocaleLowerCase())
}

const ShareLine = ({ share, isAvailable, ...props }: { share: Share, isAvailable: boolean }) => (
    <Box {...props} className={s.shareItem} mb="2" p="3">
        <Flex key={share.ticker} justify="between">
            <span>{share.name}</span>
            {isAvailable
                ? <Text color="gray">{quantToNumber(share.minPriceIncrement)}</Text>
                : <Text color="red">Недоступен</Text>}

        </Flex>
        <Text color="gray" size="1"> {share.ticker}</Text>
    </Box>
);

const SharesTriggerButton = () => (
    <Button highContrast variant="ghost" size="4" radius="full" className={`${style.button} ${s.triggerButton}`}>
        <MixerHorizontalIcon color='white' />
    </Button>
);

export const SharesPop = ({ trigger }: { trigger?: React.ReactNode }) => {
    const { shares } = useSharesFromStore();
    const schedules = useTodaysSchedules();
    const logger = useLogger({ component: 'SharesPop' });
    const [currentInstrument, setCurrentInstrument] = useCurrentInstrument();
    const schedulesByExchangeMap = useMemo<Record<string, TradingSchedule>>(() => {
        return schedules.reduce((acc, s) => ({ ...acc, [s.exchange]: s }), {})
    }, [schedules]);
    const [term, setTerm] = useState("");
    const shareFilter = useCallback((share: Share) => {
        // TODO: Вот это можно вычислять еще на ноде, когда фетчим инструменты
        // и хранить отдельно торгуемые и неторгуемые (не МОЕХ)

        const tradesBySupportedExchange = true || schedulesByExchangeMap[share.exchange];
        if (!term) return tradesBySupportedExchange;
        const fitsSearch = isContainsWithIgnoreCase(share.name, term) ||
            isContainsWithIgnoreCase(share.ticker, term) ||
            isContainsWithIgnoreCase(share.figi, term) ||
            isContainsWithIgnoreCase(share.uid, term);
        return tradesBySupportedExchange && fitsSearch;
    }, [term, schedulesByExchangeMap])
    const filteredShares = useMemo(() => (shares || []).filter(shareFilter), [shares, shareFilter])

    const onSearchChange: React.ChangeEventHandler<HTMLInputElement> = useCallback(({ target }) => {
        const term = target.value;
        setTerm(term);
    }, []);

    const isExchangeOpened = useCallback((share: Share) => {
        const exchange = share.exchange === 'MOEX_EVENING_WEEKEND' ? 'MOEX_CLOSE' : share.exchange.toUpperCase();
        const schedule = schedulesByExchangeMap[exchange];
        if (!schedule) return;

        return schedule.days[0].isTradingDay
            && schedule.days[0].endTime > new Date()
            && schedule.days[0].startTime < new Date()
    }, [filteredShares, schedulesByExchangeMap]);

    const handleInstrumentSelect = (uid: string) => {
        if (!uid) return;
        
        setCurrentInstrument(uid);
    };

    useEffect(() => {
        logger.trace('Shares popup opened');
        return () => {
            logger.trace('Shares popup closed');
        }
    }, [])

    return (
        <PopoverWindow trigger={trigger ?? <SharesTriggerButton />}>
            <Card className={s.container}>
                <SearchInput placeholder='Поиск...' onChange={onSearchChange} />
                <Container className={s.listScrollContainer}>
                    <ToggleGroup.Root type="single" orientation="vertical" onValueChange={handleInstrumentSelect} value={currentInstrument}>
                        {filteredShares.map(share => (
                            <ToggleGroup.Item key={share.uid} value={share.uid} asChild>
                                <ShareLine
                                    share={share}
                                    isAvailable={isExchangeOpened(share)}
                                />
                            </ToggleGroup.Item>
                        ))}
                    </ToggleGroup.Root>
                </Container>
            </Card>
        </PopoverWindow>
    )
}