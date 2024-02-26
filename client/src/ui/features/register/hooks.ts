import { useIpcInoke } from "../../hooks";



export const useRegister = () => useIpcInoke("REGISTER");
export const useSetAccount = () => useIpcInoke("SET_ACCOUNT");
export const useGetAccount = () => useIpcInoke("GET_ACCOUNTS");