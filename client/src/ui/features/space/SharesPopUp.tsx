import * as ScrollArea from '@radix-ui/react-scroll-area';
import { MixerHorizontalIcon } from "@radix-ui/react-icons"
import { Button, Card } from "@radix-ui/themes"
import React, { useCallback, useState } from "react"
import { PopoverWindow } from "../../components/PopoverWindow"
import style from '../../basicStyles.css';
import { Quatation, Share, TradingDay, TradingSchedule } from "../../../../grpcGW/shares";
import { quantToNumber } from '../../../node/ipcHandlers/marketdata';
import { getTodaysSchedules, useSharesFromStore } from './hooks';
import { SerarchInput } from '../../components/SearchInput';

const nanoPrecision = 1_000_000_000;
const quantToNumber = (q: Quatation | undefined): number => {
    return q ? Number(q.units + (q.nano / nanoPrecision)) : 0;
}

const isContainsWithIgnoreCase = (value: string, term: string): boolean => {
    return value.toLocaleLowerCase().includes(term.toLocaleLowerCase())
}

const ShareLine = (index: number, share: Share, isAvailable: boolean) => {
    return (
        <div key={index} style={{ marginBottom: '5px' }}>
            <div key={share.ticker} style={{ display: 'flex', justifyContent: 'space-between' }}>
                <span>{share.name}</span>
                {isAvailable
                    ? <span style={{ color: 'gray' }}>{quantToNumber(share.minPriceIncrement)}</span>
                    : <span style={{ color: 'red' }}>Недоступен</span>}

            </div>
            <div style={{ color: 'gray', fontSize: '12px' }}> {share.ticker}</div>
        </div>
    )
}

export const SharesPop = () => {
    const { sharesFromStore } = useSharesFromStore();
    const schedules = getTodaysSchedules();

    const [term, setTerm] = useState("")

    const onSearchChange = useCallback((target: EventTarget & HTMLInputElement) => {
        const term = target.value
        setTerm(term)
    }, [])

    const ShareListItem = () => {

        const isExchangeAvailable = (exchange: string): boolean => {
            console.log(schedules)
            const schedule = schedules.find((sh: TradingSchedule) => sh.exchange === exchange)
            return schedule
                ? schedule.days[0].isTradingDay
                && schedule.days[0].endTime < new Date()
                && schedule.days[0].startTime > new Date()
                : false
        }

        return (
            sharesFromStore
                .filter((share: Share) => schedules.map(schedule => schedule.exchange).includes(share.exchange))
                .filter((share: Share) => {
                    return isContainsWithIgnoreCase(share.name, term) ||
                        isContainsWithIgnoreCase(share.ticker, term) ||
                        share.uid.includes(term)
                }).map((share: Share, index) => ShareLine(index, share, isExchangeAvailable(share.exchange)))
        )
    }

    const SharesPopUpContent = () => {
        return (
            <Card style={{ padding: '15px' }} color="gray">
                <SerarchInput placeholder='Поиск...' onChange={({ target }) => {
                    onSearchChange(target)
                }} />
                <ScrollArea.Root style={{ width: '500px', height: '500px', color: 'white', overflow: 'auto' }}>
                    <ScrollArea.Viewport>
                        <ShareListItem />
                    </ScrollArea.Viewport>
                </ScrollArea.Root>
            </Card>
        )
    }

    const SharesTriggerButton = () => {
        return (
            <div style={{ marginRight: '20px' }}>
                <Button highContrast variant="ghost" size="4" radius="full" style={{ verticalAlign: 'middle' }} className={style.button}>
                    <MixerHorizontalIcon color='white' style={{ color: 'black' }} />
                </Button>
            </div>
        )
    }

    return (
        <PopoverWindow triger={SharesTriggerButton()}>
            {SharesPopUpContent()}
        </PopoverWindow>
    )
}