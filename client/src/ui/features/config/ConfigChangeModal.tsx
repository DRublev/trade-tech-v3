import type { FC } from 'react';
import React, { useCallback, useState } from 'react';
import { Modal } from '../../components/Modal';
import { ConfigForm } from './ConfigForm';
import * as ScrollArea from '@radix-ui/react-scroll-area';
import { useConfig } from './hooks';
import { useCurrentInstrument } from '../../utils/useCurrentInstrumentId';
import s from './ConfigChangeModal.css';
import { useLogger } from '../../../ui/hooks';
import { mergeObjects } from './mergeObjects';
import { InstrumentSelect } from '../space/IstrumentSelect/InstrumentSelect';
import { Box, Text } from '@radix-ui/themes';
import { useStrategy } from '../strategy';

type Props = {
    trigger: React.ReactNode;
    onSubmit?: (values: Record<string, any>) => void;
};

export const ConfigChangeModal: FC<Props> = ({ trigger, onSubmit }: Props) => {
    const [strategy] = useStrategy();
    const [instrumentId] = useCurrentInstrument();
    const { scheme, defaultValues, changeConfig } = useConfig(instrumentId, strategy);
    const [shouldClose, setShouldClose] = useState(undefined); // TODO: Костыль, надо подумать как сделать удобнее
    const logger = useLogger({ component: 'ConfigChangeModal' })

    const handleSubmit = async (rawValues: Record<string, any>) => {
        try {
            const [values] = mergeObjects(rawValues, defaultValues, scheme);

            await changeConfig(values);
            onSubmit && onSubmit(values);
            setShouldClose(true);
            setTimeout(() => setShouldClose(false));
        } catch (e) {
            logger.error('Error changing config ' + e);
            // TODO: Алерт, а лучше месседж в форму с разбором ошибки
        }
    };

    return (
        <Modal title="Настройки стратегии" close={shouldClose} trigger={trigger} actions={[]}>
            <ScrollArea.Root className={s.scrollContainer}>
                <ScrollArea.Viewport>
                    <Box mb="4">
                        <Text mb="2">Инструмент для торговли</Text>
                        <InstrumentSelect />
                    </Box>
                    <ConfigForm scheme={scheme} defaultValues={defaultValues} onSubmit={handleSubmit} />
                </ScrollArea.Viewport>
            </ScrollArea.Root>
        </Modal>
    );
};
