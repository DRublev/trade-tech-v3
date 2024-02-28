import * as ScrollArea from '@radix-ui/react-scroll-area';
import { MixerHorizontalIcon } from "@radix-ui/react-icons"
import { Button, Card, Flex } from "@radix-ui/themes"
import React, { useCallback, useEffect, useState } from "react"
import { PopoverWindow } from "../../components/PopoverWindow"
import style from '../../basicStyles.css';
import storage from "../../../node/Storage";
import { Share } from "../../../../grpcGW/shares";


const useSharesFromStore = () => {
    // const getShares = useGetSharesFromStore();
    const [sharesFromStore, setShares] = useState([]);
    const [isLoading, setIsLoading] = useState(false);

    const load = useCallback(async () => {
        try {
            if (isLoading) return;
            setIsLoading(true);
            const response: any = await window.ipc.invoke('GET_SHARES_FROM_STORE');
            setShares(response.shares)
        } catch (e) {
            console.log("get shares error: ", e);
            setShares([]);
        } finally {
            setIsLoading(false);
        }
    }, []);

    useEffect(() => {
        load();
    }, [])

    return { sharesFromStore, isLoading }
};

export const SharesPop = () => {
    const { sharesFromStore, isLoading } = useSharesFromStore();
    const SharesPopUpContent = () => {

        return (
            <Card color="gray">
                <ScrollArea.Root style={{ padding: '15px', width: '500px', height: '500px', color: 'white', overflow: 'auto' }}>
                    <ScrollArea.Viewport>
                        {sharesFromStore.map((share: Share) =>
                            <div key={share.ticker} style={{ display: 'flex', justifyContent: 'space-between' }}>
                                <span>{share.name}</span>
                                <span style={{ color: 'gray' }}>{share.minPriceIncrement?.nano}</span>
                            </div>)}
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