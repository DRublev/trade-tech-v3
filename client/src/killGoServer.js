// Будет вызван в child_process

const { exec } = require('child_process');

module.exports = function (port, protocol = 'tcp') {
    port = Number.parseInt(port);

    if (!port) {
        return Promise.reject('No port specified to kill');
    }

    if (process.platform === 'win32') {
        return exec('netstat -nao', (_, stdout) => {
            console.log('14 killGoServer', stdout);
            if (!stdout) return null

            const lines = stdout.split('\n')
            // The second white-space delimited column of netstat output is the local port,
            // which is the only port we care about.
            // The regex here will match only the local port column of the output
            const lineWithLocalPortRegEx = new RegExp(`^ *${protocol.toUpperCase()} *[^ ]*:${port}`, 'gm')
            const linesWithLocalPort = lines.filter(line => line.match(lineWithLocalPortRegEx))

            const pids = linesWithLocalPort.reduce((acc, line) => {
                const match = line.match(/(\d*)\w*(\n|$)/gm)
                return match && match[0] && !acc.includes(match[0]) ? acc.concat(match[0]) : acc
            }, [])

            return exec(`TaskKill /F /PID ${pids.join(' /PID ')}`)
        })
    }

    return exec(`kill -9 $(lsof -ti:${port})`)
};
