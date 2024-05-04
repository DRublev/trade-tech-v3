import { useMemo } from "react";
import { useAppDispatch, useAppSelector } from "../../store";
import { useSharesFromStore } from "../features/space/hooks";
import { setCurrentInstrument } from "../features/space/spaceSlice";
import type { Share } from "../../node/grpc/contracts/shares";


export const useCurrentInstrument = (): [string, (c: string) => void, Share | undefined] => {
    const dispatch = useAppDispatch();
    const { shares } = useSharesFromStore();
    const instrumentId = useAppSelector(s => s.space.currentInstrument);
    const fullInstrumentInfo = useMemo(() => (shares || []).find(s => s.uid === instrumentId), [shares, instrumentId]);

    const set = (candidate: string) => {
        if (!candidate) throw new Error('candidate is required');

        dispatch(setCurrentInstrument(candidate));
    };

    return [instrumentId, set, fullInstrumentInfo];
};