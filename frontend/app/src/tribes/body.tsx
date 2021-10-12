import React, { useEffect, useState } from 'react'
import styled from 'styled-components'
import { useObserver } from 'mobx-react-lite'
import { useStores } from '../store'
import {
  EuiFormFieldset,
  EuiLoadingSpinner,
  EuiHeader,
  EuiPopover,
  EuiSelectable,
  EuiHeaderSection,
  EuiButton,
  EuiFieldSearch,
  EuiHighlight,
} from '@elastic/eui';
import Tribe from './tribe'
import { useFuse, useIsMobile, useScroll } from '../hooks'
import { Divider, SearchTextInput } from '../sphinxUI';
import { orderBy } from 'lodash'
import Tag from './tag'
import tags from './tags'
// avoid hook within callback warning by renaming hooks
const getFuse = useFuse
const getScroll = useScroll

export default function BodyComponent() {
  const { main, ui } = useStores()
  const [selected, setSelected] = useState('')
  const [tagsPop, setTagsPop] = useState(false)

  const isMobile = useIsMobile()

  const selectedTags = ui.tags.filter(t => t.checked === 'on')
  const showTagCount = selectedTags.length > 0 ? true : false

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
      if (initial && initial.unique_name === un) setSelected(initial.uuid)
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

    const button = (<EuiButton
      iconType="arrowDown"
      iconSide="right"
      size="s"
      onClick={() => {
        ui.setTags(orderBy(ui.tags, ['checked'], ['asc']))
        setTagsPop(!tagsPop)
      }}>
      {`Tags ${showTagCount ? `(${selectedTags.length})` : ''}`}
    </EuiButton>)

    return <Body id="main" onScroll={handleScroll} style={{ paddingTop: 0 }}>
      <div style={{
        width: '100%', display: 'flex',
        justifyContent: 'space-between', alignItems: 'flex-start', padding: 20,
        height: 62
      }}>
        <Label>
        </Label>

        <div style={{ display: 'flex', alignItems: 'baseline' }}>
          <EuiPopover
            panelPaddingSize="none"
            button={button}
            isOpen={tagsPop}
            closePopover={() => setTagsPop(false)}>
            <EuiSelectable
              searchable
              options={ui.tags}
              renderOption={(option, searchValue) => <div style={{ display: 'flex', alignItems: 'center', }}>
                <Tag type={option.label} iconOnly />
                <EuiHighlight search={searchValue} style={{
                  fontSize: 11, marginLeft: 5, color: tags[option.label].color
                }}>
                  {option.label}
                </EuiHighlight>
              </div>}
              listProps={{ rowHeight: 30 }} // showIcons:false
              onChange={opts => {
                console.log(opts)
                ui.setTags(opts)
              }}>
              {(list, search) => <div style={{ width: 220 }}>
                {search}
                {list}
              </div>}
            </EuiSelectable>
          </EuiPopover>

          <div style={{ width: 20 }} />

          <SearchTextInput
            name='search'
            type='search'
            small={isMobile}
            placeholder='Search'
            value={ui.searchText}
            style={{ width: 204, height: 40, background: '#111', color: '#fff', border: 'none' }}
            onChange={e => {
              console.log('handleChange', e)
              ui.setSearchText(e)
            }}
          />
        </div>
      </div>
      <Column className="main-wrap">
        {loading && <EuiLoadingSpinner size="xl" style={{ marginTop: 20 }} />}
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
    </Body >
  }
  )
}

const Body = styled.div`
  flex:1;
  height:calc(100vh - 60px);
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
  margin-top:10px;
  // max-width:900px;
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
const Label = styled.div`
            font-family: Roboto;
            font-style: normal;
            font-weight: bold;
            font-size: 26px;
            line-height: 40px;
            /* or 154% */
            
            display: flex;
            align-items: center;
            
            /* Text 2 */
            
            color: #ffffff;`