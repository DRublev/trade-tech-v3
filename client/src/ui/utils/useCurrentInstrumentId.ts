import { useEffect, useMemo } from "react";
import { useAppDispatch, useAppSelector } from "../../store";
import { useSharesFromStore } from "../features/space/hooks";
import { setCurrentInstrument, setInitiallySetCurrentInstrument } from "../features/space/spaceSlice";
import type { Share } from "../../node/grpc/contracts/shares";
import { useIpcInvoke } from "../hooks";


export const useCurrentInstrument = (): [string, (c: string) => void, Share | undefined] => {
    const dispatch = useAppDispatch();
    const { shares } = useSharesFromStore();
    const instrumentId = useAppSelector(s => s.space.currentInstrument);
    const getCurrentInstrument = useIpcInvoke('GET_CURRENT_INSTRUMENT');
    const persistCurrentInstrument = useIpcInvoke('SET_CURRENT_INSTRUMENT');
    const initiallySetInstrument = useAppSelector(s => s.space.initiallySetCurrentInstrument);
    const fullInstrumentInfo = useMemo(() => (shares || []).find(s => s.uid === instrumentId), [shares, instrumentId]);

    const set = (candidate: string) => {
        if (!candidate) throw new Error('candidate is required');

        dispatch(setCurrentInstrument(candidate));

        setTimeout(() => {
            persistCurrentInstrument({ instrumentId: candidate }).catch(console.warn);
        }, 0);
    };

    useEffect(() => {
        if (initiallySetInstrument) return;
        getCurrentInstrument({}).then((instrumentId: string) => {
            set(instrumentId);
        })
            .catch(console.warn)
            .finally(() => {
                dispatch(setInitiallySetCurrentInstrument(true));
            });
    }, []);

    return [instrumentId, set, fullInstrumentInfo];
};