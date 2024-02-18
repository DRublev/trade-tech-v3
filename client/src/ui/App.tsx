import React from 'react';
import { createRoot } from 'react-dom/client';
import { Router } from './Router';

const root = createRoot(document.body);

root.render(<React.StrictMode>
    <Router />
</React.StrictMode>);