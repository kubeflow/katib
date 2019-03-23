import React from 'react';
import ReactDOM from 'react-dom';
import App from './components/App';
import * as serviceWorker from './serviceWorker';
import CssBaseline from '@material-ui/core/CssBaseline';
import { createMuiTheme, MuiThemeProvider } from '@material-ui/core/styles';
import configureStore from './store';
import rootSaga from './sagas';

import { HashRouter as Router } from 'react-router-dom';

import { Provider } from 'react-redux';

const store = configureStore();

store.runSaga(rootSaga);

const theme = createMuiTheme({
    palette: {
        primary: {
            main: '#000',
        },
        secondary: {
            main: '#fff',
        },
    },
});


ReactDOM.render(
    <Provider store={store}>
        <Router basename="/">
            <MuiThemeProvider theme={theme}>
                <CssBaseline />
                <App />
            </MuiThemeProvider>
        </Router>
    </Provider>
    , document.getElementById('root'));

// If you want your app to work offline and load faster, you can change
// unregister() to register() below. Note this comes with some pitfalls.
// Learn more about service workers: http://bit.ly/CRA-PWA
serviceWorker.unregister();