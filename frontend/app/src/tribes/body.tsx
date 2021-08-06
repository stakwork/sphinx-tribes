import React, { useEffect, useState } from 'react'
import styled from 'styled-components'
import { useObserver } from 'mobx-react-lite'
import { useStores } from '../store'
import {
  EuiFormFieldset,
  EuiLoadingSpinner,
} from '@elastic/eui';
import Tribe from './tribe'
import { useFuse, useScroll } from '../hooks'

// avoid hook within callback warning by renaming hooks
const getFuse = useFuse
const getScroll = useScroll

export default function BodyComponent() {
  const { main, ui } = useStores()
  const [selected, setSelected] = useState('')

  function selectTribe(uuid: string, unique_name: string) {
    setSelected(uuid)
    if (unique_name && window.history.pushState) {
      window.history.pushState({}, 'Sphinx Tribes', '/t/' + unique_name);
    }
  }
  async function loadTribes() {
    let un = ''
    if (window.location.pathname.startsWith('/t/')) {
      un = window.location.pathname.substr(3)
    }
    const ts = await main.getTribes(un)
    if (un) {
      const initial = ts[0]
      if (initial.unique_name === un) setSelected(initial.uuid)
    }
  }
  useEffect(() => {
    loadTribes()
  }, [])

  return useObserver(() => {

    const loading = main.tribes.length === 0

    const tagsFilter = ui.tags.filter(t => t.checked === 'on').map(t => t.label)
    const tribes = main.tribes.map(t => {
      const matchCount = tagsFilter.reduce((a, item) => t.tags.includes(item) ? a + 1 : a, 0)
      return { ...t, matchCount }
    }).filter(t => {
      if (tagsFilter.length === 0) return true
      return t.matchCount && t.matchCount > 0
    })

    const nsfwChecked = tagsFilter.find(label => label === 'NSFW') ? true : false
    const sfwTribes = nsfwChecked ? tribes :
      tribes.filter(t => !t.tags.includes('NSFW'))

    let theTribes = getFuse(sfwTribes, ["name", "description"])
    const { n, loadingMore, handleScroll } = getScroll()
    const finalTribes = theTribes.slice(0, n)

    return <Body id="main" onScroll={handleScroll}>
      <Column className="main-wrap">
        {loading && <EuiLoadingSpinner size="xl" />}
        {!loading && <EuiFormFieldset style={{ width: '100%' }} className="container">
          <div className="row">
            {finalTribes.map(t => <Tribe {...t} key={t.uuid}
              selected={selected === t.uuid}
              select={selectTribe}
            />)}
          </div>
        </EuiFormFieldset>}
        <LoadmoreWrap show={loadingMore}>
          <EuiLoadingSpinner size="l" />
        </LoadmoreWrap>
      </Column>
    </Body>
  }
  )
}

const Body = styled.div`
  flex:1;
  height:calc(100vh - 90px);
  padding-bottom:80px;
  width:100%;
  overflow:auto;
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
interface LoadmoreWrapProps {
  show: boolean
}
const LoadmoreWrap = styled.div<LoadmoreWrapProps>`
  position:relative;
  text-align:center;
  visibility:${p => p.show ? 'visible' : 'hidden'};
`
