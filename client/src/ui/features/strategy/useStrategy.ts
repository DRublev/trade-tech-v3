import { useState } from "react";
import { SUPPORTED_STRATEGIES } from "./constants";
import type { StrategyKey } from "./types";
export const useStrategy = (): [StrategyKey, (candidate: string) => void, Record<StrategyKey, string>] => {
    const defaultStrategy = 'spread_v0';
    const [currentStrategy, setCurrentStrategy] = useState<StrategyKey>(defaultStrategy);

    const setStrategy = (candidate: StrategyKey) => {
        if (!SUPPORTED_STRATEGIES[candidate]) {
            throw new Error("Не поддерживаемая стратегия");
        }

        setCurrentStrategy(candidate);
    }

    return [currentStrategy, setStrategy, SUPPORTED_STRATEGIES];
};
