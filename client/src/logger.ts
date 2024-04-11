import pino from "pino";
import axios from 'axios';
import path from 'path';
import fs from 'fs';

const lokiOptions = {
    address: "http://79.174.80.98:3100",
};

const mappedLevels: Record<string, string> = {
    '10': 'trace',
    '20': 'debug',
    '30': 'info',
    '40': 'warn',
    '50': 'error',
    '60': 'fatal'
}

type LogEntry = {
    labels: Record<string, string>;
    log: [string, string];
};

let chunk: LogEntry[] = [];
let pendingChunk: LogEntry[] = [];
const maxChunkSize = 30;
const sendTimeout = 10 * 1000; // 10 sec
let sendTimerId: any;

// ! Важно чтобы совпадало с путем из uniqId в го

let uid: string;

// TODO: Избавиться в пользу destination в файл. А на стороне фронта и ноды сделать воркеры, которые будут отправлять логи в локи
const toLoki = async (log: string) => {
    const { level, time, msg, ...labels } = JSON.parse(log);


    const logEntry: LogEntry = {
        labels: { app: "trade-tech", level: mappedLevels[level] || 'unknown' },
        log: [`${time}000000`, Object.keys(labels).reduce((acc, l) => acc + `${l}=${labels[l]} `, '') + ' ' + msg],
    };
    if (uid) {
        logEntry.labels.uid = uid;
    }
    chunk.push(logEntry);

    /**
     * Либо у нас заполняется чанк, либо срабатывает таймер
     * Это позволит не терять логи при простое приложения, например
     */
    if (chunk.length >= maxChunkSize) {
        if (sendTimerId) clearTimeout(sendTimerId);
        send();
    }
    if (!sendTimerId) {
        sendTimerId = setTimeout(send, sendTimeout);
    }
};

const send = () => {
    pendingChunk = chunk;
    const logs = pendingChunk.reduce((acc, logEntry) => {
        acc.streams.push({
            stream: logEntry.labels,
            values: [logEntry.log],
        })
        return acc;
    }, { streams: [] });
    chunk = [];

    axios.post(lokiOptions.address + "/loki/api/v1/push", logs).then(() => {
        pendingChunk = [];
    }).catch(e => {
        console.error('Error sending log', e.status || e.data || e);
        if (pendingChunk.length >= maxChunkSize * 10) {
            const time = Date.now();
            chunk.push({
                labels: { app: "trade-tech", level: mappedLevels[50] },
                log: [`${time}000000`, `Flushing logs of size ${pendingChunk.length}`],
            });
            pendingChunk = [];
        } else {
            // Своеобразный ретрай
            chunk = chunk.concat(pendingChunk);
        }
    });
}

const logger = pino(
    {
        browser: {
            serialize: true,
            asObject: false,
            write: (m) => toLoki(JSON.stringify(m)),
        },

    },
    {
        write: (m) => toLoki(m),
    }
);
export const setUid = (candidate: string) => {
    uid = candidate;
}
export default logger;
