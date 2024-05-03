const { exec } = require('child_process');
const path = require('path');

const UNIX_SERVER_FILENAME = "app-binary-macos";
const WIN_SERVER_FILENAME = "app-binary-windows.exe";

const serverFilename = process.platform == 'win32' ? WIN_SERVER_FILENAME : UNIX_SERVER_FILENAME;
const launch = () => {
    const isPackaged = process.argv[process.argv.length - 1] == '--packaged';

    const resourcesPath = path.resolve(isPackaged ? process.resourcesPath : '../server/resources/app/go-binaries');

    const logFileName = `${new Date().toLocaleDateString()}.log`;
    const logPath = path.join(resourcesPath, 'logs', logFileName)
    const logFileCmd = isPackaged ? '' : ` >> ${logPath}`;
    let serverPath = path
        .join(resourcesPath, serverFilename);
    if (process.platform !== 'win32') {
        serverPath = serverPath
            .split('/')
            .reduce((p, part, i) => i !== 0 ? p.concat('/', part.includes(' ') ? `'${part}'` : part) : p, '')
            .trim();
    }

    console.log('Launching server by path: ', "serverPath " + serverPath);

    try {
        const p = exec(`${serverPath}${logFileCmd}`, (err, stdout, stderr) => {
            console.log('24 launchGoServer', err);
            console.log('25 launchGoServer', stdout);
            console.log('26 launchGoServer', stderr);
        });
        p.stdout.on('data', data => {
            console.log('38 launchGoServer', data);

            if (`${data}`.includes('Server listening at')) {
                process.parentPort.postMessage('OK');
            }
        });
    } catch (e) {
        console.log('27 launchGoServer', e);
        process.parentPort.postMessage(JSON.stringify(e));

    }
};

module.exports = launch();