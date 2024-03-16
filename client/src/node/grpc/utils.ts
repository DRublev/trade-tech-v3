import { ChannelCredentials, Client, credentials } from "@grpc/grpc-js";

type UnpackedCallback<T> = T extends (err: infer E, result: infer U) => void
    ? U
    : T;
type GenericFunction<TS extends any[]> = (...args: TS) => unknown;
type Promisify<T> = {
    [K in keyof T]: T[K] extends GenericFunction<infer TS>
    ? (request: TS[0]) => Promise<UnpackedCallback<TS[3]>>
    : never;
};

const DEFAULT_ADDRESS = "0.0.0.0:50051";
const DEFAULT_CREDS = credentials.createInsecure();

export const createService = <T>(
    ClientClass: typeof Client,
    address: string = DEFAULT_ADDRESS,
    creds: ChannelCredentials = DEFAULT_CREDS
): Promisify<T> => {
    const service = new ClientClass(address, creds);
    const promisified = Object.keys(service).reduce((p: T, v: string) => {
        p[v] = (request: any) =>
            new Promise((resolve, reject) => {
                p[v](request, (err: Error, result: any) => {
                    if (err) return reject(err);
                    resolve(result);
                });
            });
        return p;
    }, {}) as unknown as Promisify<T>;
    return promisified;
};
