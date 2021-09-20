import React, { useState, useEffect } from 'react'
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

  const selectedTags = ui.tags.filter(t => t.checked === 'on')
  const showTagCount = selectedTags.length > 0 ? true : false

  useEffect(() => {
    if (window.location.host === 'podcasts.sphinx.chat') {
      ui.setTags(ui.tags.map(t => {
        if (t.label === 'Podcast') return { ...t, checked: 'on' }
        return t
      }))
    }
  }, [])

  return useObserver(() => {

    return <EuiHeader id="header" >
      {/* <div className="container"> */}
      <div className="row" style={{ marginLeft: 15 }}>
        <EuiHeaderSection grow={false} className="col-xs-12 col-sm-12 col-md-6 col-lg-6">
          <img id="logo" src="/static/tribes_logo.svg" alt="Logo" />
          {/*<Title>Tribes</Title>*/}
        </EuiHeaderSection>
      </div>
      {/* </div> */}
    </EuiHeader>
  })
}

interface ImageProps {
  readonly src: string;
  readonly size: string;
};

const Image = styled.div<ImageProps>`
  background-image:url("${(p) => p.src}");
  background-position:center;
  background-size:cover;
  height:${(p) => p.size ? p.size : '88'}px;
  width:${(p) => p.size ? p.size : '88'}px;
  position:relative;
`
const Title = styled.div`
  height:50px;
  display:flex;
  margin-left:10px;
  align-items:center;
  font-size:21px;
`