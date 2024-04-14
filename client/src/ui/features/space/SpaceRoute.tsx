import React, { useState, useRef } from 'react';
import { Layout } from "../../components/Layout"
import * as Toolbar from '@radix-ui/react-toolbar';
import { Box, Card, Flex, Section } from "@radix-ui/themes";
import { ListBulletIcon, MixerHorizontalIcon, PersonIcon, PlayIcon, StopIcon } from '@radix-ui/react-icons';
import style from '../../basicStyles.css';
import Chart from "../chart";
import s from './styles.css';
import { SharesPop } from './SharesPopup/SharesPopUp';
import { useIpcInvoke, useLogger } from '../../hooks';
import { useCurrentInstrument } from '../../utils/useCurrentInstrumentId';
import { useNavigate } from 'react-router-dom';
import { ConfigChangeModal } from '../config';
import { TradeLogs } from '../tradeSessionStats/TradeLogs';
import { TradeStats } from '../tradeSessionStats/TradeStats';
import { ipcEvents } from '../../../ipcEvents';

const toolBarButtonProps = {
    className: style.button,
    style: { verticalAlign: 'middle', transform: 'scale(1.6)', marginRight: '20px' },
}

export const ControlsPanel = () => {
    const startTrade = useIpcInvoke('START_TRADE');
    const isStartedReq = useIpcInvoke<unknown, {Ok: boolean}>(ipcEvents.IS_STARTED);
    const stopTrade = useIpcInvoke('STOP_TRADE');
    const navigate = useNavigate();
    const [instrument] = useCurrentInstrument();
    const [isStarted, setIsStarted] = useState(false);
    const [isLoading, setIsLoading] = useState(false);
    const logger = useLogger({ component: 'ControlsPanel' });

    const toggleTrade = async () => {
        try {
            logger.info("Switching trade", { isStarted })
            let res: any = {};

            if (isStarted) {
                res = await stopTrade({
                    instrumentId: instrument,
                });
            } else {
                res = await startTrade({
                    instrumentId: instrument,
                });
            }
            if (res.Ok) {
                setIsStarted(!isStarted);
            }
        } catch (e) {
            logger.error(e, "Error switching trade state");
        } finally {
            setIsLoading(false);
            logger.trace("Trade switched", { isStarted });
        }
    };

    const onAccountClick = () => {
        logger.info("Going to accounts select screen")
        navigate('/register/select-account');
    };

    React.useEffect(() => {
      isStartedReq({instrumentId: instrument})
        .then(res => setIsStarted(res.Ok));
    }, [])

    return (
        <Toolbar.Root>
            <Toolbar.ToggleGroup type="single">
                <Flex align="center" justify="center" gap="2" p="3">
                    <SharesPop
                        trigger={
                            <Toolbar.Button asChild {...toolBarButtonProps}>
                                <ListBulletIcon color='white' />
                            </Toolbar.Button>
                        }
                    />
                    <Toolbar.Button value="start" asChild onClick={toggleTrade} {...toolBarButtonProps}>
                        {isStarted ? <StopIcon color={isLoading ? "grey" : undefined} /> : <PlayIcon color={isLoading ? "grey" : undefined} />}
                    </Toolbar.Button>
                    <ConfigChangeModal
                        trigger={
                            <Toolbar.Button value="change-config" asChild {...toolBarButtonProps}>
                                <MixerHorizontalIcon color="white" />
                            </Toolbar.Button>
                        }
                    />
                    <Toolbar.Button value="logout" asChild onClick={onAccountClick} {...toolBarButtonProps}>
                        <PersonIcon />
                    </Toolbar.Button>
                </Flex>
            </Toolbar.ToggleGroup>
        </Toolbar.Root>
    )
}



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