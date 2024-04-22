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

type Props = {
    trigger: React.ReactNode;
};

export const ConfigChangeModal: FC<Props> = ({ trigger }: Props) => {
    // TODO: Брать ключ стратегии из какого-то провайдера
    const strategy = "spread_v0";
    const [instrumentId] = useCurrentInstrument();
    const { api, scheme, defaultValues, changeConfig } = useConfig(instrumentId, strategy);
    const [shouldClose, setShouldClose] = useState(undefined); // TODO: Костыль, надо подумать как сделать удобнее
    const logger = useLogger({ component: 'ConfigChangeModal' })

    const onSubmit = useCallback(async (rawValues: Record<string, any>) => {
        try {
            const [values] = mergeObjects(rawValues, defaultValues, scheme);

            await changeConfig(values);
            setShouldClose(true);
            setTimeout(() => setShouldClose(false))
        } catch (e) {
            logger.error('Error changing config', e);
            // TODO: Алерт, а лучше месседж в форму с разбором ошибки
        }
    }, [scheme, api, defaultValues]);

    return (
        <Modal title="Настройки стратегии" close={shouldClose} trigger={trigger} actions={[]}>
            <ScrollArea.Root className={s.scrollContainer}>
                <ScrollArea.Viewport>
                    <ConfigForm scheme={scheme} defaultValues={defaultValues} onSubmit={onSubmit} />
                </ScrollArea.Viewport>
            </ScrollArea.Root>
        </Modal>
    );
};
