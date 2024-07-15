import { SUPPORTED_STRATEGIES } from "./constants";
import type { StrategyKey } from "./types";
import { useAppDispatch, useAppSelector } from "../../../store";
import { setStrategy } from "../space/spaceSlice";

export const useStrategy = (): [StrategyKey, (candidate: string) => void, Record<StrategyKey, string>] => {
    const dispatch = useAppDispatch();
    const currentStrategy = useAppSelector(s => s.space.strategy);

    const selectStrategy = (candidate: StrategyKey) => {
        if (!SUPPORTED_STRATEGIES[candidate]) {
            throw new Error("Не поддерживаемая стратегия");
        }

        dispatch(setStrategy(candidate))
    }

    return [currentStrategy, selectStrategy, SUPPORTED_STRATEGIES];
};
