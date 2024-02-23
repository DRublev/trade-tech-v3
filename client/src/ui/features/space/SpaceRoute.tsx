import React, { useState, useRef, useEffect } from 'react';
import { Layout } from "../../components/Layout"
import * as Toolbar from '@radix-ui/react-toolbar';
import { Button, Card, Flex } from "@radix-ui/themes";
import { PlayIcon, StopIcon } from '@radix-ui/react-icons';
import style from '../../basicStyles.css';
import Chart from "./chart/Chart";
import s from './styles.css';
import { useIpcInoke } from '../../hooks';
import { useGetCandles } from './hooks';

export const ControlsPanel = () => {
    const [isStarted, setIsStarted] = useState(false);
    const onStartClick = () => {
        setIsStarted(!isStarted)
        //future logic
    }

    return (
        <Card style={{ minWidth: '40vw', padding: 0, position: 'fixed', left: '50%', transform: 'translate(-50%, 50%)', bottom: '40px', margin: '0 auto', boxShadow: '4px 4px 8px 0px rgba(34, 60, 80, 0.2)' }}>
            <Toolbar.Root>
                <Flex align="center" justify="center" gap="4">
                    <Toolbar.ToggleGroup type="single">
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

// BBG004730N88
const useHistoricCandels = (figiOrInstrumentId: string = "BBG004730N88") => {
    const getCandles = useGetCandles();
    const [data, setData] = useState([]);
    const [isLoading, setIsLoading] = useState(false);
    const [error, setError] = useState(null);

    const now = new Date();
    const dayAgoDate = new Date(new Date(now).setDate(now.getDate() - 1))

    const getInitialCandels = async () => {
        try {
            setIsLoading(true);

            const candles = await getCandles({
                instrumentId: figiOrInstrumentId,
                interval: 1,
                start: dayAgoDate,
                end: now,
            });
            setData(candles);
            console.log("57 SpaceRoute", candles);

        } catch (e) {
            setError(e);
            console.log("47 SpaceRoute", e);
        } finally {
            setIsLoading(false);
        }
    };

    useEffect(() => {
        getInitialCandels();
    }, [figiOrInstrumentId]);

    return { data, isLoading, error };
}

export const SpaceRoute = () => {
    const chartContainer = useRef();
    const { data, isLoading } = useHistoricCandels();

    return (
        <Layout>
            <Card ref={chartContainer} className={s.chartContainer}>
                {isLoading ? 'loading candles...' : <Chart parentRef={chartContainer} data={data} />
                }
            </Card>
            <ControlsPanel />
        </Layout>
    )
}