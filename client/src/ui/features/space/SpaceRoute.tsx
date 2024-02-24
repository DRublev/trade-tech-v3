import React, { useState } from 'react';
import { Layout } from "../../components/Layout"
import * as Toolbar from '@radix-ui/react-toolbar';
import { Button, Card, Flex } from "@radix-ui/themes";
import { PlayIcon, StopIcon } from '@radix-ui/react-icons'
import style from '../../basicStyles.css'
import { useGetShares } from "../";
export const ControlsPanel = () => {
    const getShares = useGetShares();
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
                            <Button className={style.button} onClick={onStartClick} highContrast variant="ghost" size="1" radius="full" style={{ verticalAlign: 'middle' ,transform:'scale(1.6)'}}>
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
    return (
        <Layout>

            <ControlsPanel />
        </Layout>
    )
}