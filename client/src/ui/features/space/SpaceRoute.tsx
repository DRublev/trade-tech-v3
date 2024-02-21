import React from "react"
import styles from './styles.module.css';
import { Layout } from "../../components/Layout"
import * as Toolbar from '@radix-ui/react-toolbar';
import { Button, Card, Flex } from "@radix-ui/themes";
import { PlayIcon } from '@radix-ui/react-icons'

export const ControlsPanel = () => {
    return (
        <Card style={{ minWidth: '20vw', padding: 0, position: 'fixed', bottom: '40px', margin: '0 auto', boxShadow: 'var(--shadow-3)' }}>
            <Toolbar.Root>
                <Flex align="center" justify="center" gap="4">
                    <Toolbar.ToggleGroup type="single">
                        <Toolbar.ToggleItem value="start" asChild>
                            <Button className={styles.button} highContrast variant="ghost" size="4" radius="full">
                                <PlayIcon />
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