import { useState } from "react";
import { SUPPORTED_STRATEGIES } from "./constants";
import type { StrategyKey } from "./types";

export const useStrategy = (): [StrategyKey, (candidate: string) => void, Record<StrategyKey, string>] => {
    const defaultStrategy = localStorage.getItem('strategy') || 'spread_v0';
    const [currentStrategy, setCurrentStrategy] = useState<StrategyKey>(defaultStrategy as any);

    const setStrategy = (candidate: StrategyKey) => {
        if (!SUPPORTED_STRATEGIES[candidate]) {
            throw new Error("Не поддерживаемая стратегия");
        }

        localStorage.setItem('strategy', candidate);

        setCurrentStrategy(candidate);
    }

    return [currentStrategy, setStrategy, SUPPORTED_STRATEGIES];
};
