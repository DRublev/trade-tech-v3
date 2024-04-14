import { app, shell, BrowserWindow, utilityProcess } from 'electron';
// This allows TypeScript to pick up the magic constants that's auto-generated by Forge's Webpack
// plugin that tells the Electron app where to look for the Webpack-bundled app code (depending on
// whether you're running in development or production).
declare const MAIN_WINDOW_WEBPACK_ENTRY: string;
declare const MAIN_WINDOW_PRELOAD_WEBPACK_ENTRY: string;

import './node/ipcHandlers'
import { getShares } from './node/ipcHandlers/instruments';
import logger from './logger';

// Handle creating/removing shortcuts on Windows when installing/uninstalling.
if (require('electron-squirrel-startup')) {
  app.quit();
}
const instrumentBaseState = 1
const fetchSharesList = async () => {
  try {
    await getShares({ instrumentStatus: instrumentBaseState });
  }
  catch (error) { logger.error('Fetching shares list error', error) }
}
const createWindow = (): void => {
  // TODO: Подумать над тем, чтобы вынести общение с сервером (стриминговые запросы) в воркер или отдельное спрятанное окно
  // Create the browser window.
  const mainWindow = new BrowserWindow({
    height: 600,
    width: 800,
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
  mainWindow.webContents.openDevTools();

  // mainWindow.webContents.on('will-navigate', (event, url) => {
  //   event.preventDefault();
  //   shell.openExternal(url);
  // });
};


// This method will be called when Electron has finished
// initialization and is ready to create browser windows.
// Some APIs can only be used after this event occurs.
app.on('ready', () => {
  const goLaunchProcess = utilityProcess.fork('src/launchGoServer.js', [app.isPackaged ? '--packaged' : '']);
  goLaunchProcess.once('spawn', () => {
    logger.info('go server starting');
    createWindow();
  });
  fetchSharesList();
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



// In this file you can include the rest of your app's specific main process
// code. You can also put them in separate files and import them here.
