// "BBG004730RP0" /* GAZP */
// "4c466956-d2ce-4a95-abb4-17947a65f18a" TGLD
// "BBG004730ZJ9" /* VTBR */
// "BBG004PYF2N3" /* POLY */
let instrumentId = "BBG004PYF2N3" /* POLY */;

export const useCurrentInstrument = (): [string, (c: string) => void] => {
    const set = (candidate: string) => {
        if (!candidate) throw new Error('candidate is required');

        instrumentId = candidate;
    };

    return [instrumentId, set];
};