import React, {useState} from 'react'
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
  return useObserver(() => {
    const button = (<EuiButton
      style={{width:142}}
      iconType="arrowDown"
      iconSide="right"
      size="s"
      onClick={()=>{
        ui.setTags(orderBy(ui.tags, ['checked'], ['asc']))
        setTagsPop(!tagsPop)
      }}>
      {`Tags ${showTagCount?`(${selectedTags.length})`:''}`}
    </EuiButton>)
    return <EuiHeader style={{justifyContent:'space-between',alignItems:'center',maxHeight:50,height:50,minHeight:50}}>
      <EuiHeaderSection grow={false}>
        <Image src="static/icon-1024.png" size="50" />
        <Title>Tribes</Title>
      </EuiHeaderSection>

      <EuiHeaderSection side="right" style={{display:'flex',alignItems:'center'}}>
        {/* <EuiHeaderSectionItem> */}
        <EuiPopover
          panelPaddingSize="none"
          button={button}
          isOpen={tagsPop}
          closePopover={()=>setTagsPop(false)}>
          <EuiSelectable
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
            onChange={opts=> ui.setTags(opts)}>
            {(list, search) => <div style={{ width: 220 }}>
              {search}
              {list}
            </div>}
          </EuiSelectable>
        </EuiPopover>
        <div style={{margin:'0 6px'}}>
          <EuiFieldSearch
            style={{width:'40vw'}}
            placeholder="Search Tribes"
            value={ui.searchText}
            onChange={e=> ui.setSearchText(e.target.value)}
            // isClearable={this.state.isClearable}
            aria-label="search"
          />
        </div>
      </EuiHeaderSection>

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
