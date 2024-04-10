import { Bindings, ChildLoggerOptions } from 'pino';
import logger, { setUid } from '../logger';
import path from 'path';
import fs from 'fs'

type TypedBindings = Bindings & {
    controller: string;
};

const getUid = () => {
    try {
        const rootDir = path.resolve(process.cwd(), '..');
        const uidStore = path.join(rootDir, "storage", 'uniqId');    
        const data = fs.readFileSync(uidStore, 'utf8');
    
        return data;
    } catch (e) {
        return '';
    }
};

export const createLogger = (bindings?: TypedBindings, options?: ChildLoggerOptions) => {
    const uid = getUid();
    setUid(uid);
    if (bindings) {
        return logger.child(bindings, options);
    }

    return logger;
}