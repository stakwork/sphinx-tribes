import React, {useState, useEffect} from 'react'
import { useObserver } from 'mobx-react-lite'
import { useStores } from '../store'
import { orderBy } from 'lodash'
import styled from 'styled-components'

import {
  EuiHeader,
  EuiPopover,
  EuiSelectable,
  EuiHeaderSection,
  EuiButton,
  EuiFieldSearch,
  EuiHighlight,
} from '@elastic/eui';

import Tag from './tag'
import tags from './tags'

export default function Header() {
  const { main, ui } = useStores()

  const [tagsPop, setTagsPop] = useState(false)

  const selectedTags = ui.tags.filter(t=>t.checked==='on')
  const showTagCount = selectedTags.length>0?true:false

  useEffect(()=>{
    if(window.location.host==='podcasts.sphinx.chat') {
      ui.setTags(ui.tags.map(t=>{
        if(t.label==='Podcast') return {...t,checked:'on'}
        return t
      }))
    }
  }, [])

  return useObserver(() => {
    const button = (<EuiButton
      iconType="arrowDown"
      iconSide="right"
      size="s"
      onClick={()=>{
        ui.setTags(orderBy(ui.tags, ['checked'], ['asc']))
        setTagsPop(!tagsPop)
      }}>
      {`Tags ${showTagCount?`(${selectedTags.length})`:''}`}
    </EuiButton>)
    return <EuiHeader id="header" >
      <div className="container">
        <div className="row">
          <EuiHeaderSection grow={false} className="col-xs-12 col-sm-12 col-md-6 col-lg-6">
            <img id="logo" src="static/tribes_logo.svg" alt="Logo"/>
            {/*<Title>Tribes</Title>*/}
          </EuiHeaderSection>

          <EuiHeaderSection id="header-right" side="right" className="col-xs-12 col-sm-12 col-md-6 col-lg-6">
            {/* <EuiHeaderSectionItem> */}
            <div>
              <EuiFieldSearch id="search-input"
                placeholder="Search Tribes"
                value={ui.searchText}
                onChange={e=> ui.setSearchText(e.target.value)}
                // isClearable={this.state.isClearable}
                aria-label="search"
              />
            </div>
            <EuiPopover
              panelPaddingSize="none"
              button={button}
              isOpen={tagsPop}
              closePopover={()=>setTagsPop(false)}>
              <EuiSelectable className="popover-tags" 
                searchable
                options={ui.tags}
                renderOption={(option,searchValue)=><div style={{display:'flex',alignItems:'center'}}>
                  <Tag type={option.label} iconOnly />
                  <EuiHighlight search={searchValue} style={{
                    fontSize:11,marginLeft:5,color:tags[option.label].color
                  }}>
                    {option.label}
                  </EuiHighlight>
                </div>}
                listProps={{rowHeight: 30}} // showIcons:false
                onChange={opts=> {
                  console.log(opts)
                  ui.setTags(opts)
                }}>
                {(list, search) => <div style={{ width: 220 }}>
                  {search}
                  {list}
                </div>}
              </EuiSelectable>
            </EuiPopover>
          </EuiHeaderSection>
        </div>
      </div>
    </EuiHeader>
  })
}

interface ImageProps {
  readonly src: string;
  readonly size: string;
};

const Image = styled.div<ImageProps>`
  background-image:url("${(p)=> p.src}");
  background-position:center;
  background-size:cover;
  height:${(p)=> p.size?p.size:'88'}px;
  width:${(p )=> p.size?p.size:'88'}px;
  position:relative;
`
const Title = styled.div`
  height:50px;
  display:flex;
  margin-left:10px;
  align-items:center;
  font-size:21px;
`
