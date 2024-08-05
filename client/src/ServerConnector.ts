import type { Logger } from "pino";
import * as ffi from 'ffi-rs'
import path from "path";

const LIB_MAME = 'app-binary';
const UNIX_SERVER_FILENAME = `${LIB_MAME}-macos.so`;
const WIN_SERVER_FILENAME = `${LIB_MAME}-windows.dll`;
const libFilename = process.platform == 'win32' ? WIN_SERVER_FILENAME : UNIX_SERVER_FILENAME;

export class ServerConnector {
    constructor(private logger: Logger, private app: Electron.App) {
        // if (!(process.env.ENV === "PROD" || app.isPackaged)) return;
        app.on('will-quit', () => {
            this.stop();
        })
    }

    public async Start(): Promise<boolean> {
        try {
            const libDir = path.resolve(this.app.isPackaged ? process.resourcesPath : '../server');

            ffi.open({
                path: path.join(libDir, libFilename),
                library: LIB_MAME,
            });
            const res = await ffi.load({
                library: LIB_MAME,
                funcName: 'LaunchServer',
                retType: ffi.DataType.I32,
                paramsType: [],
                paramsValue: [],
                runInNewThread: false,
            });
            console.log('39 ServerConnector', res);
            return true;
        } catch (e) {
            console.log('41 ServerConnector', e);
            return false;
        }
    }

    private stop() {
        ffi.close(LIB_MAME);
    }
} 