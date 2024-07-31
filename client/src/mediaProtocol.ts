import { app, net, protocol } from "electron";
import path from "path";

/**
 * Регистрируем кастомный протокол, чтобы ходить за статикой в папку ресурсов, вместо прямого урла к файлу
 */
export const registerMediaProtocol = () => {
    protocol.registerSchemesAsPrivileged([{
        scheme: 'trademedia',

        privileges: { bypassCSP: true }
    }])
    app.whenReady().then(() => {
        protocol.handle('trademedia', (request) => {
            const baseDir = app.isPackaged ? process.resourcesPath : __dirname;
            const { pathname, hostname } = new URL(request.url)
            const resourcePath = path.join(baseDir, hostname, pathname);

            return net.fetch('file://' + resourcePath)
        })
    })
};