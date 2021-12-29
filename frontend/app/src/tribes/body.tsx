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
import { useFuse, useIsMobile, useScroll, usePageScroll } from '../hooks'
import { Divider, SearchTextInput } from '../sphinxUI';
import { orderBy } from 'lodash'
import Tag from './tag'
import tags from './tags'
import NoResults from '../people/utils/noResults'
import PageLoadSpinner from '../people/utils/pageLoadSpinner';
// avoid hook within callback warning by renaming hooks
// const getFuse = useFuse
// const getScroll = useScroll

const getPageScroll = usePageScroll
let debounceValue: any = []

export default function BodyComponent() {
  const { main, ui } = useStores()
  const [selected, setSelected] = useState('')
  const [tagsPop, setTagsPop] = useState(false)
  const [tagOptions, setTagOptions] = useState(ui.tags)
  const [loading, setLoading] = useState(true)

  const { tribesPageNumber } = ui

  const isMobile = useIsMobile()

  const selectedTags = tagOptions.filter(t => t.checked === 'on')
  const showTagCount = selectedTags.length > 0 ? true : false


  function selectTribe(uuid: string, unique_name: string) {
    setSelected(uuid)
    if (unique_name && window.history.pushState) {
      window.history.pushState({}, 'Sphinx Tribes', '/t/' + unique_name);
    }
  }

  async function loadMore(direction) {
    let currentPage = tribesPageNumber
    let newPage = currentPage + direction
    if (newPage < 1) newPage = 1

    try {
      await main.getTribes({ page: newPage })
    } catch (e) {
      console.log(e)
    }
  }

  // do search update
  useEffect(() => {
    (async () => {

      console.log('refresh list')
      // reset page will replace all results, this is good for a new search!
      await main.getTribes({ page: 1, resetPage: true })

      // do deeplink
      let deeplinkUn = ''
      if (window.location.pathname.startsWith('/t/')) {
        deeplinkUn = window.location.pathname.substr(3)
      }
      if (deeplinkUn) {
        let t = await main.getTribeByUn(deeplinkUn)
        setSelected(t.uuid)
        window.history.pushState({}, 'Sphinx Tribes', '/t');
      }

      setLoading(false)
    })()
  }, [ui.searchText, ui.tags])

  function doDelayedValueUpdate() {
    ui.setTags(debounceValue)
  }

  return useObserver(() => {

    const tribes = main.tribes

    const loadForwardFunc = () => loadMore(1)
    const loadBackwardFunc = () => loadMore(-1)
    const { loadingTop, loadingBottom, handleScroll } = getPageScroll(loadForwardFunc, loadBackwardFunc)

    if (loading) {
      return <Body style={{ justifyContent: 'center', alignItems: 'center', background: '#212529' }}>
        <EuiLoadingSpinner size="xl" />
      </Body>
    }

    const button = (<EuiButton
      iconType="arrowDown"
      iconSide="right"
      size="s"
      onClick={() => {
        setTagsPop(!tagsPop)
      }}>
      {`Tags ${showTagCount ? `(${selectedTags.length})` : ''}`}
    </EuiButton>)

    return <Body id="main" onScroll={tagsPop ? () => console.log('scroll') : handleScroll} style={{ paddingTop: 0 }}>
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
              options={tagOptions}
              renderOption={(option, searchValue) => <div style={{ display: 'flex', alignItems: 'center', }}>
                <Tag type={option.label} iconOnly />
                <EuiHighlight search={searchValue} style={{
                  fontSize: 11, marginLeft: 5, ///color: ui.tags[option.label].color
                }}>
                  {option.label}
                </EuiHighlight>
              </div>}
              listProps={{ rowHeight: 30 }} // showIcons:false
              onChange={opts => {
                console.log(opts)
                setTagOptions(opts)
                debounceValue = opts
                debounce(doDelayedValueUpdate, 800)

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
        <PageLoadSpinner show={loadingTop} />
        <EuiFormFieldset style={{ width: '100%', paddingBottom: 0 }} className="container">
          <div className="row">
            {tribes.length ? tribes.map(t => <Tribe {...t} key={t.uuid}
              selected={selected === t.uuid}
              select={selectTribe}
            />) : <NoResults />}
          </div>
        </EuiFormFieldset>
        <PageLoadSpinner noAnimate show={loadingBottom} />
      </Column>
    </Body >
  }
  )
}


let inDebounce
function debounce(func, delay) {
  clearTimeout(inDebounce)
  inDebounce = setTimeout(() => {
    func()
  }, delay)
}

const Body = styled.div`
  flex:1;
  height:calc(100vh - 60px);
  // padding-bottom:80px;
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