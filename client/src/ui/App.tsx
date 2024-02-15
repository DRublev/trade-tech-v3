import React from 'react';
import { createRoot } from 'react-dom/client';
import { Router } from './Router';

const root = createRoot(document.body);
window.addEventListener('load', () => {
    window.ipc.invoke('TEST_HELLO', { asda: 'asd' })
})
root.render(<React.StrictMode>
    <Router />
</React.StrictMode>);