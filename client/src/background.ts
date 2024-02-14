import { app, BrowserWindow } from 'electron';
app.on('ready', () => {
    console.log('Hi from background process')
});
