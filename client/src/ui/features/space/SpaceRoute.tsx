import React, { useState, useRef } from 'react';
import { Layout } from "../../components/Layout"
import * as Toolbar from '@radix-ui/react-toolbar';
import { Button, Card, Flex } from "@radix-ui/themes";
import { PlayIcon, StopIcon } from '@radix-ui/react-icons';
import style from '../../basicStyles.css';
import Chart from "./chart";
import s from './styles.css';
import { SharesPop } from './SharesPopUp';
import { useIpcInoke } from '../../hooks';
import { useCurrentInstrumentId } from './hooks';

export const ControlsPanel = () => {
    const startTrade = useIpcInoke('START_TRADE');
    const [isStarted, setIsStarted] = useState(false);
    const [instrument] = useCurrentInstrumentId();

    const onStartClick = async () => {
        setIsStarted(!isStarted)
        //future logic
        try {

            await startTrade({
                instrumentId: instrument,
            });
        } catch (e) {
            console.log('24 SpaceRoute', e);
        }
    }

    return (
        <Card style={{ minWidth: '40vw', padding: 0, position: 'fixed', left: '50%', transform: 'translate(-50%, 50%)', bottom: '40px', margin: '0 auto', boxShadow: '4px 4px 8px 0px rgba(34, 60, 80, 0.2)' }}>
            <Toolbar.Root>
                <Flex align="center" justify="center" gap="4">
                    <Toolbar.ToggleGroup type="single">
                        <SharesPop />
                        <Toolbar.ToggleItem value="start" asChild>
                            <Button className={style.button} onClick={onStartClick} highContrast variant="ghost" size="1" radius="full" style={{ verticalAlign: 'middle', transform: 'scale(1.6)' }}>
                                {isStarted ? <StopIcon /> : <PlayIcon />}
                            </Button>
                        </Toolbar.ToggleItem>
                    </Toolbar.ToggleGroup>
                </Flex>
            </Toolbar.Root>
        </Card>
    )
}



export const SpaceRoute = () => {
    const chartContainer = useRef();

    return (
        <Layout>
            <Card ref={chartContainer} className={s.chartContainer}>
                <Chart containerRef={chartContainer} />
            </Card>
            <ControlsPanel />
        </Layout>
    )
}