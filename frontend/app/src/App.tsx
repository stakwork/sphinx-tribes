import React, { useEffect } from 'react'
import './App.css'
import '@material/react-material-icon/dist/material-icon.css';
import "@fontsource/roboto";
import Header from './tribes/header'
import Body from './tribes/body'
import PeopleHeader from './people/header'
import MobilePeopleHeader from './people/mobile/header'
import PeopleBody from './people/body'
import MobilePeopleBody from './people/mobile/body'
import { colors } from './colors'
import {
  BrowserRouter as Router,
  Switch,
  Route,
  Link
} from "react-router-dom";
import { useIsMobile } from './hooks';
import { MeInfo } from './store/ui';
import api from './api';
import { useStores } from './store';

function App() {
  const mode = getMode()
  const c = colors['light']
  const isMobile = useIsMobile()
  const { main, ui } = useStores()

  async function testChallenge(chal: string) {
    try {

      const me: MeInfo = await api.get(`poll/${chal}`)
      if (me && me.pubkey) {
        ui.setMeInfo(me)
        ui.setEditMe(true)
      }
    } catch (e) {
      console.log(e)
    }
  }

  useEffect(() => {
    try {
      var urlObject = new URL(window.location.href);
      var params = urlObject.searchParams;
      const chal = params.get('challenge')
      if (chal) {
        testChallenge(chal)
      }
    } catch (e) { }
  }, [])

  return <Router>
    {
      // people
      mode === Mode.PEOPLE ? <div className="app" style={{ background: c.background }}>
        {/* {isMobile ? */}
        <MobilePeopleHeader />
        {/* : <PeopleHeader />} */}
        {/* {isMobile ? */}
        <MobilePeopleBody />
        {/* : <PeopleBody />} */}
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
