import { useIpcInvoke } from "../../hooks";



export const useRegister = () => useIpcInvoke("REGISTER");
export const useSetAccount = () => useIpcInvoke("SET_ACCOUNT");
export const useGetAccount = () => useIpcInvoke("GET_ACCOUNTS");