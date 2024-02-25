import React, { useState, useRef } from 'react';
import { Layout } from "../../components/Layout"
import * as Toolbar from '@radix-ui/react-toolbar';
import { Button, Card, Flex } from "@radix-ui/themes";
import { PlayIcon, StopIcon } from '@radix-ui/react-icons';
import style from '../../basicStyles.css';
import Chart from "./chart/Chart";
import s from './styles.css';
import { useCandles } from './hooks';
import { useGetShares } from './hooks';

export const ControlsPanel = () => {

    const getShares = useGetShares();
    const [isStarted, setIsStarted] = useState(false);
    const onStartClick = async () => {

        setIsStarted(!isStarted)
        const response: any = await getShares({ instrumentStatus: 1 });

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



export const SpaceRoute = () => {
    const chartContainer = useRef();
    // TODO: Прокидывать id выбранного инструмента
    // const { data, isLoading } = useCandles();

    return (
        <Layout>
            <Card ref={chartContainer} className={s.chartContainer}>
                {/* TODO: Запилить лоадер для графика
                    Может быть прокидывать isLoading внутрь и разруливать внутри Chart
                */}
                {/* {isLoading ? 'loading candles...' : <Chart parentRef={chartContainer} data={data} /> */}
                {/* } */}
            </Card>
            <ControlsPanel />
        </Layout>
    )
}