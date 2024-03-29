import { useIpcInvoke } from "../../hooks";
import { ConfigScheme } from "./types";

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
                type: 'number',
                min: 0,
                htmlType: 'number',
            },
            {
                name: 'MaxSharesToHold',
                required: true,
                label: 'Максимально лотов',
                placeholder: 'штук',
                type: 'number',
                min: 1,
                step: 1,
                htmlType: 'number',
            },
            {
                name: 'MinProfit',
                required: true,
                label: 'Минимальный профит со сделки',
                placeholder: '',
                type: 'number',
                min: 0,
                step: 0.01,
                htmlType: 'number',
            },
            {
                name: 'StopLossAfter',
                label: 'Стоп-лосс',
                placeholder: 'цена покупки - стоп-лосс = цена продажи',
                type: 'number',
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
    api: ReturnType<typeof useConfigIpc>;
    scheme: ReturnType<typeof useConfigScheme>;
    defaultValues: Record<string, any>;
}
export const useConfig: UseConfigHook = (instrumentId: string, strategy: string) => {
    const api = useConfigIpc();
    const scheme = useConfigScheme(instrumentId, strategy);

    const defaultValues = {
        Balance: 450,
        MaxSharesToHold: 1,
        MinProfit: 0.34,
        StopLossAfter: 1,
    };

    return { api, scheme, defaultValues };
}