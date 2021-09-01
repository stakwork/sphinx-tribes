import React, { useEffect, useState } from 'react'
import './App.css'
import '@material/react-material-icon/dist/material-icon.css';
import "@fontsource/roboto";
import Header from './tribes/header'
import Body from './tribes/body'
import PeopleHeader from './people/mobile/header'
import PeopleBody from './people/mobile/body'
import { colors } from './colors'
import {
  BrowserRouter as Router,
  Switch,
  Route,
} from "react-router-dom";
import TokenRefresh from './people/utils/tokenRefresh';
import BotsBody from './bots/body';

function App() {
  const mode = getMode()
  const c = colors['light']

  return <Router>
    {
      // people
      mode === Mode.PEOPLE ? <div className="app" style={{ background: c.background }}>
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
}
const hosts: { [k: string]: Mode } = {
  "localhost:3000": Mode.TRIBES,
  "localhost:13000": Mode.TRIBES,
  "localhost:23000": Mode.TRIBES,
  "tribes.sphinx.chat": Mode.TRIBES,
  "tribes-test.sphinx.chat": Mode.TRIBES,
  "localhost:13007": Mode.PEOPLE,
  "localhost:23007": Mode.PEOPLE,
  "localhost:3007": Mode.PEOPLE,
  "people.sphinx.chat": Mode.PEOPLE,
  "people-test.sphinx.chat": Mode.PEOPLE,
};

function getMode(): Mode {

  const host = window.location.host;

  return hosts[host] || Mode.TRIBES;
}

export default App;
