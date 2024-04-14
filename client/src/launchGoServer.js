const { exec } = require('child_process');
const path = require('path');


const UNIX_SERVER_FILENAME = "server";
const WIN_SERVER_FILENAME = "server.exe";
const serverFilename = process.platform == 'win32' ? WIN_SERVER_FILENAME : UNIX_SERVER_FILENAME;
const launch = () => {
    const isPackaged = process.argv[process.argv.length - 1] == '--packaged';
    const logFileCmd = isPackaged ? '' : ' >> logs';
    const serverPath = path.join(isPackaged ? process.resourcesPath : '..', serverFilename);
    console.log('Launching server by path: ', serverPath);
    exec(`ENV=PROD ./${serverFilename}${logFileCmd}`, (err, stdout, stderr) => {
        console.log('Launch server out:', err, stderr, stdout);
    })
};

module.exports = launch();