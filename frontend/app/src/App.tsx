import React, { useEffect } from 'react';
/* eslint-disable func-style */
import '@material/react-material-icon/dist/material-icon.css';
import { BrowserRouter as Router } from 'react-router-dom';
import './App.css';
import { ModeDispatcher } from './config/ModeDispatcher';
import { Pages } from './pages';
import { mainStore } from './store/main';
import { WithModalStore } from 'store/modals';

let exchangeRateInterval: any = null;

function App() {
  // get usd/sat exchange rate every 100 seconds
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
    <WithModalStore>
      <Router>
        <ModeDispatcher>{(mode) => <Pages mode={mode} />}</ModeDispatcher>
      </Router>
    </WithModalStore>
  );
}

export default App;
