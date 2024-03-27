import type { FC } from 'react';
import React from 'react';
import { Modal } from '../../components/Modal';
import { Form } from './Form';
import * as ScrollArea from '@radix-ui/react-scroll-area';
import { useConfigScheme } from './hooks';


type Props = {
    trigger: React.ReactNode;
};

export const ConfigChangeModal: FC<Props> = ({ trigger }: Props) => {
    const instrument = useCurrent
    const { scheme } = useConfigScheme();

    return (
        <Modal trigger={trigger} title="Настройки стратегии">
            <ScrollArea.Root style={{ width: '500px', height: '500px', color: 'white', overflow: 'auto' }}>
                <ScrollArea.Viewport>
                    <Form scheme={scheme} />
                </ScrollArea.Viewport>
            </ScrollArea.Root>
        </Modal>
    );
};
