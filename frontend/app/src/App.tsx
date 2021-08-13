import React from 'react'
import './App.css'
import '@material/react-material-icon/dist/material-icon.css';
import "@fontsource/roboto";
import Header from './tribes/header'
import Body from './tribes/body'
import PeopleHeader from './people/header'
import PeopleBody from './people/body'
import { colors } from './colors'

function App() {
  const mode = getMode()
  const c = colors['light']

  if (mode === Mode.PEOPLE) {
    return <div className="app" style={{ background: c.background }}>
      <PeopleHeader />
      <PeopleBody />
    </div>
  }
  return (
    <div className="app" style={{ background: c.background }}>
      <Header />
      <Body />
    </div>
  );
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
