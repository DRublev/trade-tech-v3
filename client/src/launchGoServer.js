const { exec } = require('child_process');
const path = require('path');
const killOnPort = require('./killGoServer');

const UNIX_SERVER_FILENAME = "app-binary-macos";
const WIN_SERVER_FILENAME = "app-binary-windows.exe";

const serverFilename = process.platform == 'win32' ? WIN_SERVER_FILENAME : UNIX_SERVER_FILENAME;

const PORT = 50051;

const killPreviousServerInstance = () => new Promise((resolve, reject) => {
    const killP = killOnPort(PORT);

    killP.on('exit', resolve);
    killP.on('error', reject);
})

const launch = () => {
    const isPackaged = process.argv[process.argv.length - 1] == '--packaged';
    const resourcesPath = path.resolve(isPackaged ? process.resourcesPath : '../server/resources/app/go-binaries');
    let serverPath = path
        .join(resourcesPath, serverFilename);
    if (process.platform !== 'win32') {
        serverPath = serverPath
            .split('/')
            .reduce((p, part, i) => i !== 0 ? p.concat('/', part.includes(' ') ? `'${part}'` : part) : p, '')
            .trim();
    }

    killPreviousServerInstance().then(() => {
        const cmd = `${serverPath} --port ${PORT}`
        console.log('Launching server by path: ', cmd);
        try {
            const p = exec(cmd);
            p.stdout.on('data', data => {
                if (`${data}`.includes('Server listening at')) {
                    process.parentPort.postMessage('OK');
                }
            });
        } catch (e) {
            process.parentPort.postMessage(JSON.stringify(e));
        }
    }).catch(process.parentPort.postMessage)
};

module.exports = launch();