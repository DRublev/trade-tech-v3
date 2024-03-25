import React, { useState, useRef } from 'react';
import { Layout } from "../../components/Layout"
import * as Toolbar from '@radix-ui/react-toolbar';
import { Card, Flex } from "@radix-ui/themes";
import { MixerHorizontalIcon, PersonIcon, PlayIcon, StopIcon } from '@radix-ui/react-icons';
import style from '../../basicStyles.css';
import Chart from "./chart";
import s from './styles.css';
import { SharesPop } from './SharesPopUp';
import { useIpcInoke } from '../../hooks';
import { useNavigate } from 'react-router-dom';

const toolBarButtonProps = {
  className: style.button,
  style: { verticalAlign: 'middle', transform: 'scale(1.6)', marginRight: '20px' },
}

export const ControlsPanel = () => {
    const startTrade = useIpcInoke('START_TRADE');
    const [isStarted, setIsStarted] = useState(false);
    const navigate = useNavigate();

    const onStartClick = async () => {
        setIsStarted(!isStarted)
        //future logic
        try {

            await startTrade({});
        } catch (e) {
            console.log('24 SpaceRoute', e);
        }
    }

    const onAccountClick = () => {
      navigate('/register/select-account');
    }

    return (
        <Toolbar.Root>
            <Flex align="center" justify="center" gap="4">
                <Toolbar.ToggleGroup type="single">
                    <SharesPop 
                      trigger={
                        <Toolbar.Button asChild {...toolBarButtonProps}>
                          <MixerHorizontalIcon color='white' style={{ color: 'black' }} />
                        </Toolbar.Button>
                      } 
                    />
                    <Toolbar.Button value="start" asChild onClick={onStartClick} {...toolBarButtonProps}>
                        {isStarted ? <StopIcon /> : <PlayIcon />}
                    </Toolbar.Button>
                    <Toolbar.Button value="start" asChild onClick={onAccountClick} {...toolBarButtonProps}>
                        <PersonIcon />
                    </Toolbar.Button>
                </Toolbar.ToggleGroup>
            </Flex>
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
            <Card>
              <ControlsPanel />
            </Card>
        </Layout>
    )
}