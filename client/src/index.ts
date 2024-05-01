import { app, BrowserWindow, utilityProcess, dialog, MessageBoxOptions, ipcMain, autoUpdater, shell } from 'electron';
import os from 'os';
import { getShares } from './node/ipcHandlers/instruments';
import logger from './logger';
import { createUpdateYml } from './createUpdateYml';
import { registerMediaProtocol } from './mediaProtocol';

import './node/ipcHandlers';

// This allows TypeScript to pick up the magic constants that's auto-generated by Forge's Webpack
// plugin that tells the Electron app where to look for the Webpack-bundled app code (depending on
// whether you're running in development or production).
declare const MAIN_WINDOW_WEBPACK_ENTRY: string;
declare const MAIN_WINDOW_PRELOAD_WEBPACK_ENTRY: string;


// Handle creating/removing shortcuts on Windows when installing/uninstalling.
if (require('electron-squirrel-startup')) {
  app.quit();
}

const handleStartupEvent = function () {
  if (process.platform !== 'win32') {
    return false;
  }

  const squirrelCommand = process.argv[1];
  switch (squirrelCommand) {
    case '--squirrel-install':
    case '--squirrel-firstrun':
    case '--squirrel-updated':

      // Optionally do things such as:
      //
      // - Install desktop and start menu shortcuts
      // - Add your .exe to the PATH
      // - Write to the registry for things like file associations and
      //   explorer context menus

      // Always quit when done
      app.quit();

      return true;
    case '--squirrel-uninstall':
      // Undo anything you did in the --squirrel-install and
      // --squirrel-updated handlers

      // Always quit when done
      app.quit();

      return true;
    case '--squirrel-obsolete':
      // This is called on the outgoing version of your app before
      // we update to the new version - it's the opposite of
      // --squirrel-updated
      app.quit();
      return true;
  }
};

const runGoServer = () => {
  let serverProcess: Electron.UtilityProcess;

  app.on('will-quit', () => {
    if (!serverProcess) return;
    serverProcess.kill();
  });

  return new Promise((resolve) => {
    if (process.env.ENV === 'PROD' || app.isPackaged) {
      let scriptPath = 'src/launchGoServer.js';
      if (app.isPackaged) {
        scriptPath = process.resourcesPath + '/launchGoServer.js'
      }
      scriptPath = scriptPath
        .split('/')
        .reduce((p, part, i) => i !== 0 ? p.concat('/', part.includes(' ') ? `'${part}'` : part) : part, '')
        .trim();

      const serverProcess = utilityProcess.fork(scriptPath, [app.isPackaged ? '--packaged' : '']);
      serverProcess.once('spawn', () => {
        logger.info('go server starting');
      });
      serverProcess.on('message', m => {
        if (m === 'OK') {
          return resolve(true);
        }
      });

      serverProcess.on('exit', (code) => {
        logger.info(`go server exited with code ${code}`);
      });
    }

    return resolve(true);
  });
};

const instrumentBaseState = 1;
const fetchSharesList = async () => {
  try {
    await getShares({ instrumentStatus: instrumentBaseState });
  }
  catch (error) { logger.error('Fetching shares list error ' + error); }
};

createUpdateYml();
registerMediaProtocol();

const platform = os.platform() + '_' + os.arch();
const version = app.getVersion();
const server = 'http://79.174.80.98:5000';
const url = `${server}/update/${platform}/${version}`;

autoUpdater.setFeedURL({
  url,
  headers: {
    Accept: 'application/vnd.github+json',
    Authorization: 'Bearer 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',
    'X-GitHub-Api-Version': '2022-11-28'
  }
});


let mainWindow: BrowserWindow;
const createWindow = (): void => {
  // TODO: Подумать над тем, чтобы вынести общение с сервером (стриминговые запросы) в воркер или отдельное спрятанное окно
  // Create the browser window.
  mainWindow = new BrowserWindow({
    height: 600,
    width: 800,
    title: `Trade Tech ${app.getVersion()}`,
    webPreferences: {
      preload: MAIN_WINDOW_PRELOAD_WEBPACK_ENTRY,
      nodeIntegrationInWorker: true,
      // Use pluginOptions.nodeIntegration, leave this alone
      // See nklayman.github.io/vue-cli-plugin-electron-builder/guide/security.html#node-integration for more info
      nodeIntegration: !!process.env.ELECTRON_NODE_INTEGRATION,
      contextIsolation: !process.env.ELECTRON_NODE_INTEGRATION,
    },
  });

  // and load the index.html of the app.
  mainWindow.loadURL(MAIN_WINDOW_WEBPACK_ENTRY);

  // Open the DevTools.
  if (!app.isPackaged) {
    mainWindow.webContents.openDevTools();
  }

  mainWindow.webContents.on('will-navigate', (event, url) => {
    event.preventDefault();
    shell.openExternal(url);
  });
};



// This method will be called when Electron has finished
// initialization and is ready to create browser windows.
// Some APIs can only be used after this event occurs.
app.on('ready', async () => {
  if (handleStartupEvent()) return;

  runGoServer().then(hasLaunched => {
    if (!hasLaunched) return;

    createWindow();
    fetchSharesList();
  });

  const onWindowsOnlyIfPacked = process.platform == 'win32' && app.isPackaged;
  if (onWindowsOnlyIfPacked) {
    autoUpdater.checkForUpdates();
  }
});

// Quit when all windows are closed, except on macOS. There, it's common
// for applications and their menu bar to stay active until the user quits
// explicitly with Cmd + Q.
app.on('window-all-closed', () => {
  if (process.platform !== 'darwin') {
    app.quit();
  }
});

app.on('activate', () => {
  // On OS X it's common to re-create a window in the app when the
  // dock icon is clicked and there are no other windows open.
  if (BrowserWindow.getAllWindows().length === 0) {
    createWindow();
  }
});

ipcMain.handle('RESIZE', (e, req) => {
  // eslint-disable-next-line prefer-const
  let { width, height } = req;

  if (mainWindow.webContents.devToolsWebContents) {
    width += 300;
  }

  mainWindow.setSize(width, height);
  mainWindow.center();
})

autoUpdater.on('update-downloaded', (event, releaseNotes, releaseName, releaseDate, updateURL) => {
  const dialogOpts: MessageBoxOptions = {
    type: 'info',
    buttons: ['Перезапустить', 'Позже'],
    title: 'Приложение обновлено',
    message: Array.isArray(releaseNotes) ? [...releaseNotes].map(String).join('\n') : releaseNotes,
    detail: 'Новая версия приложения была загружена. Перезапустить приложение для применения обновления?',
  }

  dialog.showMessageBox(dialogOpts).then((returnValue) => {
    if (returnValue.response === 0) {
      setTimeout(() => {
        autoUpdater.quitAndInstall();
      }, 1_000);
    }
  });
});

autoUpdater.on('error', (error) => {
  logger.error('Ошибка обновления' + error);
  dialog.showErrorBox('Ошибка обновления', error.message);
})

// In this file you can include the rest of your app's specific main process
// code. You can also put them in separate files and import them here.
