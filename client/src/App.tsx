import React from 'react';
import { createRoot } from 'react-dom/client';

const root = createRoot(document.body);
window.addEventListener('load', () => {
    window.ipc.invoke('TEST_HELLO', { asda: 'asd' })
})
root.render(<h2>Hello from React!</h2>);