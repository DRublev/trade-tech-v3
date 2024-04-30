import pino from "pino";
import axios from 'axios';

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

let uid: string;

// TODO: Избавиться в пользу destination в файл. А на стороне фронта и ноды сделать воркеры, которые будут отправлять логи в локи
const toLoki = async (log: Record<string, any> | string) => {
    const logObj = typeof log === "string" ? JSON.parse(log) : log;
    const { level, time, msg, pid, ...labels } = logObj;

    const logEntry: LogEntry = {
        labels: { ...labels, pid: `${pid}`, app: "client", level: mappedLevels[level] || 'unknown', message: msg },
        log: [`${time}000000`, `${msg}`],
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
        console.error('Error sending log', e.status || (e.response && e.response.data) || e);
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
            serialize: false,
            asObject: true,
            write: (m) => toLoki(m),
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
