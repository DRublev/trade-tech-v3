import { Client } from "@grpc/grpc-js";
import { DEFAULT_ADDRESS, DEFAULT_CREDS } from "./constants";

type UnpackedCallback<T> = T extends (err: infer E, result: infer U) => void
    ? U
    : T;
type GenericFunction<TS extends any[]> = (...args: TS) => unknown;

type Promisify<T> = {
    [K in keyof T]: T[K] extends GenericFunction<infer TS>
    ? (request: TS[0]) => Promise<UnpackedCallback<TS[3]>>
    : never;
};

export const createService = <T extends Client>(
    ClientClass: typeof Client,
): Promisify<T> => {
    const service = new ClientClass(DEFAULT_ADDRESS, DEFAULT_CREDS)
    const promisified: Promisify<T> = {} as any;
    for (const v in service) {
        promisified[v as keyof T] = ((request: unknown) =>
            new Promise((resolve, reject) => {
                (service[v as keyof Client] as any)(request, (err: Error, result: any) => {
                    if (err) return reject(err);
                    resolve(result);
                });
            })) as any;
    }
    return promisified;
};
