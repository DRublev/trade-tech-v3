import React, { useRef } from 'react';
import { Layout } from "../../components/Layout"
import { Card, Flex, Section } from "@radix-ui/themes";
import Chart from "../chart";
import s from './styles.css';
import { TradeLogs } from '../tradeSessionStats/TradeLogs';
import { TradeStats } from '../tradeSessionStats/TradeStats';
import { ControlsPanel } from './ControlsPanel/ControlsPanel';


export const SpaceRoute = () => {
    const chartContainer = useRef();

    return (
        <Layout>
            <Card ref={chartContainer} className={s.chartContainer}>
                <Chart containerRef={chartContainer} />
            </Card>
            <Flex p="0" gap="4" justify="between">
                <Section width="auto" my="2" p="1" className={s.tradeLogsContainer}>
                    <TradeStats />
                </Section>
                <Section width="auto" my="2" p="1" className={s.tradeLogsContainer}>
                    <TradeLogs />
                </Section>
            </Flex>
            <Card className={s.controlsContainer}>
                <ControlsPanel />
            </Card>
        </Layout>
    )
}