import { ipcEvents } from "../../../ipcEvents";
import { useIpcInvoke } from "../../hooks";



export const useRegister = () => useIpcInvoke(ipcEvents.REGISTER);
export const usePruneTokens = () => useIpcInvoke(ipcEvents.PRUNE_TOKENS);
export const useSetAccount = () => useIpcInvoke(ipcEvents.SET_ACCOUNT);
export const useGetAccount = () => useIpcInvoke(ipcEvents.GET_ACCOUNTS);