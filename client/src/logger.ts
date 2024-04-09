import pino from "pino";
import axios from 'axios'

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

const toLoki = async (log: string) => {
    const { level, time, msg, ...labels } = JSON.parse(log)

    axios.post(lokiOptions.address + "/loki/api/v1/push", {
        streams: [{
            stream: { app: "trade-tech", level: mappedLevels[level] || 'unknown', ...Object.keys(labels).reduce((acc, l) => ({ ...acc, [l]: `${labels[l]}` }), {}) }, values: [
                [`${time}000000`, msg]
            ]
        }]
    }).catch(e => {
        console.error('Error sendning log', e.status || e.data);
    });
};

const logger = pino(
    {
        browser: {
            serialize: true,
            asObject: false,
        },
    },
    {
        write: (m) => toLoki(m),
    }
);

export default logger;
