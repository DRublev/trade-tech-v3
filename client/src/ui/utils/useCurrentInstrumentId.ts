import { useAppDispatch, useAppSelector } from "../../store";
import { setCurrentInstrument } from "../features/space/spaceSlice";


export const useCurrentInstrument = (): [string, (c: string) => void] => {
    const dispatch = useAppDispatch();
    const instrumentId = useAppSelector(s => s.space.currentInstrument);

    const set = (candidate: string) => {
        if (!candidate) throw new Error('candidate is required');

        dispatch(setCurrentInstrument(candidate));
    };

    return [instrumentId, set];
};