import * as ScrollArea from '@radix-ui/react-scroll-area';
import { MixerHorizontalIcon } from "@radix-ui/react-icons"
import { Button, Card } from "@radix-ui/themes"
import React, { useState } from "react"
import { PopoverWindow } from "../../components/PopoverWindow"
import style from '../../basicStyles.css';
import { Quatation, Share } from "../../../../grpcGW/shares";
import { quantToNumber } from '../../../node/ipcHandlers/marketdata';
import { useSharesFromStore } from './hooks';
import { SerarchInput } from '../../components/SearchInput';

const nanoPrecision = 1_000_000_000;
const quantToNumber = (q: Quatation | undefined): number => {
    return q ? Number(q.units + (q.nano / nanoPrecision)) : 0;
}

const isContainsWithIgnoreCase = (value: string, term: string): boolean => {
    return value.toLocaleLowerCase().includes(term) ||
        value.toUpperCase().includes(term) ||
        value.includes(term)
}
export const SharesPop = () => {
    const { sharesFromStore } = useSharesFromStore();

    const [term, setTerm] = useState("")

    const onSearchChange = (target: EventTarget & HTMLInputElement) => {
        const term = target.value
        setTerm(term)
    }

    const mapShares = () => {
        return (
            sharesFromStore.filter((share: Share) => {
                return isContainsWithIgnoreCase(share.name, term) ||
                    isContainsWithIgnoreCase(share.ticker, term) ||
                    share.uid.includes(term)
            }).map((share: Share) =>
                <div style={{ marginBottom: '5px' }}>
                    <div key={share.ticker} style={{ display: 'flex', justifyContent: 'space-between' }}>
                        <span>{share.name}</span>
                        <span style={{ color: 'gray' }}>{quantToNumber(share.minPriceIncrement)}</span>
                    </div>
                    <div style={{ color: 'gray', fontSize: '12px' }}> {share.ticker}</div>
                </div>

            )
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
                        {mapShares()}
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