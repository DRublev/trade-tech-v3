import React from 'react';
import { createRoot } from 'react-dom/client';
import { Router } from './Router';
import { store } from '../store';
import { Provider } from 'react-redux';

const root = createRoot(document.querySelector('#root') || document.body);

root.render(<React.Fragment>
    <Provider store={store}>
        <Router />
    </Provider>
</React.Fragment>);