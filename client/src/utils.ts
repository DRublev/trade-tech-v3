import { Quant } from "./node/ipcHandlers/types";

const nanoPrecision = 1_000_000_000;
export const quantToNumber = (q: Quant): number => {
    return Number(q.units + (q.nano / nanoPrecision));
}