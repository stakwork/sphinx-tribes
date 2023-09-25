import React, { useCallback, useEffect } from 'react';
/* eslint-disable func-style */
import '@material/react-material-icon/dist/material-icon.css';
import history from 'config/history';
import { withProviders } from 'providers';
import { Router } from 'react-router-dom';
import { uiStore } from 'store/ui';
import './App.css';
import { ThemeProvider, createTheme } from '@mui/system';
import { ModeDispatcher } from './config/ModeDispatcher';
import { Pages } from './pages';
import { mainStore } from './store/main';

let exchangeRateInterval: any = null;

const theme = createTheme({
  spacing: 8
});

function App() {
  const getUserOrganizations = useCallback(async () => {
    if (uiStore.meInfo && uiStore.meInfo?.tribe_jwt) {
      await mainStore.getUserOrganizations();
    }
  }, []);

  useEffect(() => {
    getUserOrganizations();
  }, [getUserOrganizations]);

  useEffect(() => {
    // get usd/sat exchange rate every 100 second;
    mainStore.getUsdToSatsExchangeRate();

    exchangeRateInterval = setInterval(() => {
      mainStore.getUsdToSatsExchangeRate();
    }, 100000);

    return function cleanup() {
      clearInterval(exchangeRateInterval);
    };
  }, []);

  return (
    <ThemeProvider theme={theme}>
      <Router history={history}>
        <ModeDispatcher>{(mode: any) => <Pages mode={mode} />}</ModeDispatcher>
      </Router>
    </ThemeProvider>
  );
}

export default withProviders(App);
