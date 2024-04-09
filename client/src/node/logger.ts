import { Bindings, ChildLoggerOptions } from 'pino';
import logger from '../logger';

type TypedBindings = Bindings & {
    controller: string;
};

export const createLogger = (bindings?: TypedBindings, options?: ChildLoggerOptions) => {
    if (bindings) {
        return logger.child(bindings, options);
    }

    return logger;
}