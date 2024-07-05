import type { FC } from 'react';
import React from 'react';
import type { StrategyKey } from './types';
import { Select } from '@radix-ui/themes';
import { useStrategy } from './useStrategy';

type Props = {
    disabled: boolean;
    onChange?: (strategy: StrategyKey) => void;
};

export const StrategySelector: FC<Props> = ({ disabled, onChange }) => {
    const [strategy, setStrategy, allStrategies] = useStrategy();
    const options = Object.entries(allStrategies).map(([key, name]) => ({ value: key, name }));

    const handleStrategyChange = (candidate: string) => {
        setStrategy(candidate);
        onChange && onChange(candidate as StrategyKey);
    };

    return <>
        <Select.Root defaultValue={strategy} onValueChange={handleStrategyChange} disabled={disabled}>
            <Select.Trigger />
            <Select.Content>
                {options.map((option) => (
                    <Select.Item key={option.value} value={option.value}>{option.name}</Select.Item>
                ))}
            </Select.Content>
        </Select.Root>
    </>;
};
