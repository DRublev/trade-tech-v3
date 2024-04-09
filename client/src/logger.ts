import pino from 'pino';

const send = (level, logObj) => {
console.log('4 logger', level, logObj);

};

const logger = pino({
    browser: {
        serialize: true,
        asObject: false,
        transmit: {
            send
        }
    }
});

export default logger;
