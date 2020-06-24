import React from 'react';
import ReactDOM from 'react-dom';
import App from './components/App';
import KubeflowDashboard from './components/KubeflowDashboard';
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
  colors: {
    created: '#33adff',
    running: '#0911f6',
    restarting: '#1eb9af',
    succeeded: '#02970a',
    failed: '#e62e00',
    killed: '#ff5c33',
  },
  typography: {
    fontFamily:
      'open sans,-apple-system,BlinkMacSystemFont,segoe ui,Roboto,helvetica neue,Arial,sans-serif,apple color emoji,segoe ui emoji,segoe ui symbol',
  },
});

ReactDOM.render(
  <Provider store={store}>
    <Router basename="/">
      <MuiThemeProvider theme={theme}>
        <KubeflowDashboard />
        <CssBaseline />
        <App />
      </MuiThemeProvider>
    </Router>
  </Provider>,
  document.getElementById('root'),
);
