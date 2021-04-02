import React from 'react'
import {useObserver} from 'mobx-react-lite'
import './App.css'
import Header from './components/header'
import Body from './components/body'

function App() {
  return useObserver(()=>
    <div className="app">
      <Header />
      <Body />
    </div>
  )
}


export default App
