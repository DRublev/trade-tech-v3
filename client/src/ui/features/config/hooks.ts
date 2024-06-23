import { useEffect, useState } from "react";
import { useIpcInvoke, useLogger } from "../../hooks";
import { ConfigScheme } from "./types";
import { DEFAULT_SCHEME, schemes } from "./schemes";
import type { StrategyKey } from "../strategy/types";


type UseConfigSchemeHook = (insrumentId: string, strategy?: StrategyKey) => ConfigScheme;

const useConfigScheme: UseConfigSchemeHook = (instrumentId, strategy) => {
    return (strategy && schemes[strategy]) || DEFAULT_SCHEME;
};


const useConfigIpc = () => ({
    change: useIpcInvoke('CHANGE_STRATEGY_CONFIG'), get: useIpcInvoke('GET_STRATEGY_CONFIG')
});


type UseConfigHook = (instrumentId: string, strategy: string) => {
    scheme: ReturnType<typeof useConfigScheme>;
    defaultValues: Record<string, any>;
    changeConfig: (values: Record<string, string>) => Promise<void>
}

export const useConfig: UseConfigHook = (instrumentId: string, strategy: StrategyKey) => {
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