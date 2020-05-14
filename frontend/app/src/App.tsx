import React, {useEffect  } from 'react'
import {useObserver} from 'mobx-react-lite'
import {useStores} from './store'
import './App.css'
import Header from './components/header'
import Body from './components/body'

function App() {
  const { main } = useStores()
  useEffect(()=>{
    main.getTribes()
  },[])
  return useObserver(()=>
    <div className="app">
      <Header />
      <Body />
    </div>
  )
}


export default App
