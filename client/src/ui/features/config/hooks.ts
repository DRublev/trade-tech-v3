import { useEffect, useState } from "react";
import { useIpcInvoke, useLogger } from "../../hooks";
import { ConfigScheme } from "./types";

type FieldTypes = ConfigScheme['fields'][number]['type'];
export const ConfigFieldTypes: Record<FieldTypes, FieldTypes> = {
    number: 'number',
    money: 'money',
}

type UseConfigSchemeHook = (insrumentId: string, strategy?: string) => ConfigScheme;

const useConfigScheme: UseConfigSchemeHook = (instrumentId, strategy) => {
    // TODO: Брать схему с бека, как будет надо вообще
    const scheme: ConfigScheme = {
        fields: [
            {
                name: 'Balance',
                required: true,
                label: 'Доступный для торговли баланс',
                placeholder: 'рублей',
                type: ConfigFieldTypes.number,
                min: 0,
                htmlType: 'number',
            },
            {
                name: 'MaxSharesToHold',
                required: true,
                label: 'Максимально лотов',
                placeholder: 'штук',
                type: ConfigFieldTypes.number,
                min: 1,
                step: 1,
                htmlType: 'number',
            },
            {
                name: 'MinProfit',
                required: true,
                label: 'Минимальный профит со сделки',
                placeholder: '',
                type: ConfigFieldTypes.number,
                min: 0,
                step: 0.01,
                htmlType: 'number',
            },
            {
                name: 'StopLossAfter',
                label: 'Стоп-лосс',
                placeholder: 'цена покупки - стоп-лосс = цена продажи',
                type: ConfigFieldTypes.number,
                min: 0,
                step: 0.01,
                htmlType: 'number',
            },
        ]
    };
    return scheme;
};


const useConfigIpc = () => ({
    change: useIpcInvoke('CHANGE_STRATEGY_CONFIG'), get: useIpcInvoke('GET_STRATEGY_CONFIG')
});


type UseConfigHook = (instrumentId: string, strategy: string) => {
    scheme: ReturnType<typeof useConfigScheme>;
    defaultValues: Record<string, any>;
    changeConfig: (values: Record<string, string>) => Promise<void>
}

export const useConfig: UseConfigHook = (instrumentId: string, strategy: string) => {
    const api = useConfigIpc();
    const scheme = useConfigScheme(instrumentId, strategy);
    const [defaultValues, setDefaultValues] = useState({});
    const logger = useLogger({ component: 'useConfig' })

    const fetchInitialValues = async () => {
        try {
            const cfg = await api.get({ instrumentId, strategy });

            setDefaultValues(cfg);
        } catch (e) {
            logger.error("Error fetching default config " + e);
        }
    }

    const changeConfig = async (values: Record<string, string>) => {
        await api.change({ instrumentId, strategy, values });
        await fetchInitialValues();
    }

    useEffect(() => {
        fetchInitialValues();
    }, [instrumentId, strategy]);

    return { scheme, defaultValues, changeConfig };
}