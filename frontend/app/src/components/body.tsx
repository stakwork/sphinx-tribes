import React, {useState} from 'react'
import styled from 'styled-components'
import { useObserver } from 'mobx-react-lite'
import { useStores } from '../store'
import {
  EuiFormFieldset,
  EuiLoadingSpinner,
} from '@elastic/eui';
import Tribe from './tribe'
import Fuse from 'fuse.js'

const fuseOptions = {
  keys: ['name','description'],
  shouldSort: true,
  // matchAllTokens: true,
  includeMatches: true,
  threshold: 0.35,
  location: 0,
  distance: 100,
  maxPatternLength: 32,
  minMatchCharLength: 1,
};

export default function BodyComponent() {
  const { main, ui } = useStores()
  const [selected, setSelected] = useState('')
  return useObserver(() => {
    const loading = main.tribes.length===0

    const tagsFilter = ui.tags.filter(t=>t.checked==='on').map(t=>t.label)
    const tribes = main.tribes.map(t=>{
      const matchCount = tagsFilter.reduce((a,item)=>t.tags.includes(item)?a+1:a, 0)
      return {...t,matchCount}
    }).filter(t=>{
      if(tagsFilter.length===0) return true
      return t.matchCount&&t.matchCount>0
    })

    let theTribes = tribes
    if(ui.searchText){
      var fuse = new Fuse(tribes, fuseOptions)
      const res = fuse.search(ui.searchText)
      theTribes = res.map(r=>r.item)
    }

    return <Body id="main">
      <Column className="main-wrap">
        {loading && <EuiLoadingSpinner size="xl" />}
        {!loading && <EuiFormFieldset style={{width:'100%'}} className="container">
          <div className="row">
            {theTribes.map(t=> <Tribe {...t} key={t.uuid}
              selected={selected===t.uuid}
              select={setSelected}
            />)}
          </div>
        </EuiFormFieldset>}
      </Column>
    </Body>
  }
)}

const Body = styled.div`
  flex:1;
  height:calc(100vh - 50px);
  padding-bottom:80px;
  width:100%;
  overflow:scroll;
  background:#272c4b;
  display:flex;
  flex-direction:column;
  align-items:center;
`
const Column = styled.div`
  display:flex;
  flex-direction:column;
  align-items:center;
  max-width:900px;
  width:100%;
`
