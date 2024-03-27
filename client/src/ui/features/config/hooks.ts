import { ConfigScheme } from "./types";

type UseConfigSchemeHook = (insrumentId: string, strategy: string) => {
    scheme: ConfigScheme;
}

export const useConfigScheme: UseConfigSchemeHook = () => {
    const scheme = { fields: [] };

    const update = (value) => {

    };

    return { scheme }
}