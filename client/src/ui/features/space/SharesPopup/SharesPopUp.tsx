import * as ScrollArea from '@radix-ui/react-scroll-area';
import { MixerHorizontalIcon } from "@radix-ui/react-icons";
import { Box, Button, Card, Flex, Section, Text } from "@radix-ui/themes";
import React, { useCallback, useMemo, useState } from "react";
import { PopoverWindow } from "../../../components/PopoverWindow";
import style from '../../../basicStyles.css';
import { Quatation, Share, TradingSchedule } from "../../../../node/grpc/contracts/shares";
import { useTodaysSchedules, useSharesFromStore } from '../hooks';
import { SearchInput } from '../../../components/SearchInput';
import s from './styles.css';

const nanoPrecision = 1_000_000_000;
const quantToNumber = (q: Quatation | undefined): number => {
    return q ? Number(q.units + (q.nano / nanoPrecision)) : 0;
}

const isContainsWithIgnoreCase = (value: string, term: string): boolean => {
    return value.toLocaleLowerCase().includes(term.toLocaleLowerCase())
}

const ShareLine = ({ share, isAvailable }: { share: Share, isAvailable: boolean }) => {
    return (
        <Box mb="1">
            <Flex key={share.ticker} justify="between">
                <span>{share.name}</span>
                {isAvailable
                    ? <Text color="gray">{quantToNumber(share.minPriceIncrement)}</Text>
                    : <Text color="red">Недоступен</Text>}

            </Flex>
            <Text color="gray" size="1"> {share.ticker}</Text>
        </Box>
    )
}

const SharesTriggerButton = () => (
    <Button highContrast variant="ghost" size="4" radius="full" className={`${style.button} ${s.triggerButton}`}>
        <MixerHorizontalIcon color='white' />
    </Button>
);

export const SharesPop = ({ trigger }: { trigger?: React.ReactNode }) => {
    const { sharesFromStore } = useSharesFromStore();
    const schedules = useTodaysSchedules();
    const schedulesByExchangeMap = useMemo<Record<string, TradingSchedule>>(() => {
        return schedules.reduce((acc, s) => ({ ...acc, [s.exchange]: s }), {})
    }, [schedules]);
    const [term, setTerm] = useState("");
    const shareFilter = useCallback((share: Share) => {
        // TODO: Вот это можно вычислять еще на ноде, когда фетчим инструменты
        // и хранить отдельно торгуемые и неторгуемые (не МОЕХ)

        const tradesBySupportedExchange = true || schedulesByExchangeMap[share.exchange];
        const fitsSearch = !term || isContainsWithIgnoreCase(share.name, term) ||
            isContainsWithIgnoreCase(share.ticker, term) ||
            share.uid.includes(term);
        return tradesBySupportedExchange && fitsSearch;
    }, [term])
    const filteredShares = useMemo(() => sharesFromStore.filter(shareFilter), [sharesFromStore, shareFilter])

    const onSearchChange: React.ChangeEventHandler<HTMLInputElement> = useCallback(({ target }) => {
        const term = target.value;
        setTerm(term);
    }, []);

    const isExchangeOpened = useCallback((share: Share) => {
        if (!schedulesByExchangeMap[share.exchange]) return;
        const schedule = schedulesByExchangeMap[share.exchange];
        return schedule.days[0].isTradingDay
            && schedule.days[0].endTime < new Date()
            && schedule.days[0].startTime > new Date()
    }, [filteredShares, schedulesByExchangeMap]);


    return (
        <PopoverWindow trigger={trigger ?? <SharesTriggerButton />}>
            <Card className={s.container}>
                <SearchInput placeholder='Поиск...' onChange={onSearchChange} />
                <ScrollArea.Root className={s.listScrollContainer}>
                    <ScrollArea.Viewport>
                        {filteredShares.map(share => (
                            <ShareLine
                                key={share.uid}
                                share={share}
                                isAvailable={isExchangeOpened(share.exchange)}
                            />
                        ))}
                    </ScrollArea.Viewport>
                </ScrollArea.Root>
            </Card>
        </PopoverWindow>
    )
}