import type { FC } from 'react';
import React, { useCallback, useState } from 'react';
import { Modal } from '../../components/Modal';
import { ConfigForm } from './ConfigForm';
import * as ScrollArea from '@radix-ui/react-scroll-area';
import { ConfigFieldTypes, useConfig } from './hooks';
import { useCurrentInstrument } from '../../utils/useCurrentInstrumentId';
import s from './ConfigChangeModal.css';

type Props = {
    trigger: React.ReactNode;
};

export const ConfigChangeModal: FC<Props> = ({ trigger }: Props) => {
    // TODO: Брать ключ стратегии из какого-то провайдера
    const strategy = "spread_v0";
    const [instrumentId] = useCurrentInstrument();
    const { api, scheme, defaultValues } = useConfig(instrumentId, strategy);
    const [shouldClose, setShouldClose] = useState(undefined); // TODO: Костыль, надо подумать как сделать удобнее

    const onSubmit = useCallback(async (rawValues: Record<string, any>) => {
        try {
            let changedFields = 0;
            const values: Record<string, any> = {};
            for (const fieldKey in rawValues) {
                if (!(fieldKey in defaultValues) || rawValues[fieldKey] != defaultValues[fieldKey]) {
                    changedFields++;
                }
                const field = scheme.fields.find(f => f.name === fieldKey);
                if (!field) continue;
                if (field.type === ConfigFieldTypes.number || field.type === ConfigFieldTypes.money) {
                    values[fieldKey] = Number(rawValues[fieldKey]);
                    continue;
                }
                values[fieldKey] = rawValues[fieldKey];
            }
            if (!changedFields) {
                setShouldClose(true);
                setTimeout(() => setShouldClose(false))
                return;
            }

            await api.change({ instrumentId, strategy, values });
            setShouldClose(true);
            setTimeout(() => setShouldClose(false))
        } catch (e) {
            console.error('22 ConfigChangeModal', e);
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
