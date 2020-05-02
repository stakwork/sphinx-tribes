import React, {useState} from 'react'
import { useObserver } from 'mobx-react-lite'
import { useStores } from '../store'
import { orderBy } from 'lodash'
import styled from 'styled-components'

import {
  EuiHeader,
  EuiHeaderBreadcrumbs,
  EuiPopover,
  EuiPopoverTitle,
  EuiSelectable,
  EuiHeaderSection,
  EuiHeaderSectionItem,
  EuiHeaderSectionItemButton,
  EuiHeaderLogo,
  EuiButton,
  EuiIcon,
  EuiFieldSearch,
  EuiComboBox,
  EuiComboBoxOptionProps,
} from '@elastic/eui';

// import HeaderAppMenu from './header_app_menu';
// import HeaderUserMenu from './header_user_menu';
// import HeaderSpacesMenu from './header_spaces_menu';

export default function Header() {
  const { main } = useStores()
  const [text, setText] = useState<string>('')
  const [selectedTags, setSelectedTags] = useState<any[]>([
    {label:'hi'},{label:'lo'}
  ])
  const [tagsPop, setTagsPop] = useState(false)

  const button = (
    <EuiButton
      iconType="arrowDown"
      iconSide="right"
      size="s"
      onClick={()=>{
        setSelectedTags(orderBy(selectedTags, ['checked'], ['asc']))
        setTagsPop(!tagsPop)
      }}>
      Tags
    </EuiButton>
  );

  return useObserver(() =>
    <EuiHeader style={{justifyContent:'space-between',alignItems:'center',maxHeight:50,height:50,minHeight:50}}>
      <EuiHeaderSection grow={false}>
        {/* <EuiHeaderSectionItem border="right">
          <Image src="static/icon-1024.png" size="50" />
          Sphinx Tribes
        </EuiHeaderSectionItem>
        <EuiHeaderSectionItem border="right">
        </EuiHeaderSectionItem> */}
        <Image src="static/icon-1024.png" size="50" />
        <Title>Sphinx Tribes</Title>
      </EuiHeaderSection>

      <EuiHeaderSection side="right" style={{display:'flex',alignItems:'center'}}>
        {/* <EuiHeaderSectionItem> */}
        <EuiPopover
          panelPaddingSize="none"
          button={button}
          isOpen={tagsPop}
          closePopover={()=>setTagsPop(false)}>
          <EuiSelectable
            options={selectedTags}
            onChange={opts=>{
              console.log(opts)
              setSelectedTags(opts)
            }}>
            {(list, search) => (
              <div style={{ width: 240 }}>
                {list}
              </div>
            )}
          </EuiSelectable>
        </EuiPopover>
        <div style={{margin:'0 6px'}}>
          <EuiFieldSearch
            style={{width:'50vw'}}
            placeholder="Search Tribes"
            value={text}
            onChange={e=> setText(e.target.value)}
            // isClearable={this.state.isClearable}
            aria-label="Use aria labels when no actual label is in use"
          />
        </div>
      </EuiHeaderSection>

    </EuiHeader>
  )
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
