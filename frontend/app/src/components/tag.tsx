import React from 'react'
import styled from 'styled-components'
import tags from './tags'

export default function T(props:any){
  const {type} = props
  if(!tags[type]) return <></>
  const Icon = tags[type].icon
  const color = tags[type].color
  return <Wrap>
    <IconWrap style={{borderColor:color,background:color+'22'}}>
      <Icon height="10" width="10" />
    </IconWrap>
    {!props.iconOnly && <Name style={{color}}>{type}</Name>}
  </Wrap>
}

const Wrap = styled.div`
  display:flex;
  align-items:center;
  margin-left:9px;
`
const Name = styled.span`
  font-size:10px;
  margin-left:3px;
`
const IconWrap = styled.div`
  width:16px;
  height:16px;
  border-width:1px;
  border-style:solid;
  border-radius:3px;
  display:flex;
  align-items:center;
  justify-content:center;
`