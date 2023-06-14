import React, { useEffect } from 'react';
/* eslint-disable func-style */
import '@material/react-material-icon/dist/material-icon.css';
import { Router } from 'react-router-dom';
import { WithStores } from './store';
import './App.css';
import { ModeDispatcher } from './config/ModeDispatcher';
import { Pages } from './pages';
import { mainStore } from './store/main';
import history from 'config/history';

let exchangeRateInterval: any = null;

function App() {
  // get usd/sat exchange rate every 100 second;

  useEffect(() => {
    mainStore.getUsdToSatsExchangeRate();

    exchangeRateInterval = setInterval(() => {
      mainStore.getUsdToSatsExchangeRate();
    }, 100000);

    return function cleanup() {
      clearInterval(exchangeRateInterval);
    };
  }, []);

  return (
    <WithStores>
      <Router history={history}>
        <ModeDispatcher>{(mode: any) => <Pages mode={mode} />}</ModeDispatcher>
      </Router>
    </WithStores>
  );
}

export default App;
