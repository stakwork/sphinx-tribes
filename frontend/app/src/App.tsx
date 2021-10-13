import React, { useEffect, useState } from 'react'
import './App.css'
import '@material/react-material-icon/dist/material-icon.css';
import "@fontsource/roboto";
import Header from './tribes/header'
import Body from './tribes/body'
import PeopleHeader from './people/main/header'
import PeopleBody from './people/main/body'
import { colors } from './colors'
import {
  BrowserRouter as Router,
  Switch,
  Route,
} from "react-router-dom";
import TokenRefresh from './people/utils/tokenRefresh';
import BotsBody from './bots/body';
import { mainStore } from './store/main';

let exchangeRateInterval: any = null

function App() {
  const mode = getMode()
  const c = colors['light']

  // get usd/sat exchange rate every 100 seconds
  useEffect(() => {
    mainStore.getUsdToSatsExchangeRate()

    exchangeRateInterval = setInterval(() => {
      mainStore.getUsdToSatsExchangeRate()
    }, 100000)

    return function cleanup() {
      clearInterval(exchangeRateInterval)
    }
  }, [])

  return <Router>
    {
      // people/community
      mode === Mode.COMMUNITY ? <div className="app" style={{ background: c.background }}>
        <PeopleHeader />
        <TokenRefresh />
        <Switch>
          <Route path="/p/">
            <PeopleBody />
          </Route>
          <Route path="/t/">
            <Body />
          </Route>
          <Route path="/b/">
            <BotsBody />
          </Route>
        </Switch>

        {/* for global toasts */}

      </div>
        // tribes
        : <div className="app" style={{ background: c.background }}>
          <Header />
          <Body />
        </div>
    }

  </Router>


}

enum Mode {
  TRIBES = "tribes",
  PEOPLE = "people",
  COMMUNITY = "community",
}

const hosts: { [k: string]: Mode } = {
  "localhost:3000": Mode.TRIBES,
  "localhost:13000": Mode.TRIBES,
  "localhost:23000": Mode.TRIBES,
  "tribes.sphinx.chat": Mode.TRIBES,
  "tribes-test.sphinx.chat": Mode.TRIBES,
  "localhost:13007": Mode.COMMUNITY,
  "localhost:23007": Mode.COMMUNITY,
  "localhost:3007": Mode.COMMUNITY,
  "people.sphinx.chat": Mode.COMMUNITY,
  "people-test.sphinx.chat": Mode.COMMUNITY,
  "community-test.sphinx.chat": Mode.COMMUNITY,
  "community.sphinx.chat": Mode.COMMUNITY,
};

function getMode(): Mode {

  const host = window.location.host;

  return hosts[host] || Mode.TRIBES;
}

export default App;
