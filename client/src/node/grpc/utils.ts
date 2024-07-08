import { Client, type CallOptions, type ClientReadableStream } from "@grpc/grpc-js";
import { DEFAULT_ADDRESS, DEFAULT_CREDS } from "./constants";

type UnpackedCallback<T> = T extends (err: infer E, result: infer U) => void
    ? U
    : T;
type GenericFunction<TS extends any[]> = (...args: TS) => unknown;

type StreamFunction<Req, Ret> = (req: Req, options?: Partial<CallOptions>) => ClientReadableStream<Ret>

type OnStreamEvent<T> = //((eN: 'data', list: (e: T) => void) => any) |
    ((...args: Parameters<ClientReadableStream<T>['on']>) => ReturnType<ClientReadableStream<T>['on']>)

type Stream<S extends ClientReadableStream<T>, T> = Omit<S, 'on'> & {
    on: OnStreamEvent<T> | ClientReadableStream<T>['on']
}

type Promisify<T> = {
    [K in keyof T]:
    // Если это стрим, то не оборачиваем в промис
    T[K] extends StreamFunction<infer SA, infer SR>
    ? StreamFunction<SA, SR>
    : T[K] extends GenericFunction<infer TS>
    ? (request: TS[0]) => Promise<UnpackedCallback<TS[3]>>
    : never;
};

/**
 * Оборачивает все унарные запросы в промисы
 * @param ClientClass Класс gRPC сервиса, окторый будет пропатчен
 * @returns Новый инстанс сервиса
 */

export const createService = <T extends Client>(
    ClientClass: typeof Client,
): Promisify<T> => {
    const service = new ClientClass(DEFAULT_ADDRESS, DEFAULT_CREDS)
    const promisified: Promisify<T> = {} as any;
    for (const v in service) {
        if ((service[v as keyof Client] as any)['responseStream']) {
            promisified[v as keyof T] = (service[v as keyof Client] as any).bind(service);
            continue;
        }

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
