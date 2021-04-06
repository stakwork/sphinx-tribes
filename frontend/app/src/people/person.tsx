import React, {useRef, useState} from 'react'
import { QRCode } from 'react-qr-svg';
import styled from 'styled-components'
import {EuiCheckableCard} from '@elastic/eui';
import moment from 'moment'

export default function Person({id,img,tags,description,selected,select,created,owner_alias,unique_name}:any){
  return <EuiCheckableCard className="col-md-6 col-lg-6 ml-2 mb-2"
    id={id+''}
    label={owner_alias}
    name={owner_alias}
    value={id+''}
    checked={selected}
    onChange={() => select(id, unique_name)}>
    <Content onClick={() => select(selected?'':id,unique_name)} style={{
      height:selected?'auto':100
    }} selected={selected}>
      <Left>
        <Row className="item-cont">
          <div className="placeholder-img-tribe"></div>
          <Img src={img} />
          <Left style={{padding:'0 0 0 20px', maxWidth:'calc(100% - 100px)'}}>
            <Row style={selected?{flexDirection:'column',alignItems:'flex-start'}:{}}>
              <Title className="tribe-title">
                {owner_alias}
              </Title>
            </Row>
            <Description oneLine={selected?false:true} style={{minHeight:20}}>
              {description}
            </Description>
          </Left>
        </Row>
        <div className="expand-part" style={selected ? { opacity: 1} : { opacity: 0}}>

          <div className="colapse-button"><img src="/static/keyboard_arrow_up-black-18dp.svg" alt="" /></div>
        </div>
      </Left>
    </Content>
  </EuiCheckableCard>
}
interface ContentProps {
  selected: boolean;
}
const Content = styled.div<ContentProps>`
  cursor:pointer;
  display:flex;
  justify-content:space-between;
  max-width:100%;
  & h3{
    color:#fff;
  }
  &:hover h3{
    color:white;
  }
  ${p=> p.selected?`
    & h5{
      color:#cacaca;
    }
  `:`
    & h5{
      color:#aaa;
    }
    &:hover h5{
      color:#bebebe;
    }
  `}
`
const QRWrap = styled.div`
  background:white;
  padding:5px;
`
const Left = styled.div`
  height:100%;
  max-width:100%;
  display:flex;
  flex-direction:column;
  flex:1;
`
const Row=styled.div`
  display:flex;
  align-items:center;
`
const Title=styled.h3`
  margin-right:12px;
  font-size:22px;
  white-space:nowrap;
  overflow:hidden;
  text-overflow:ellipsis;
  max-width:100%;
  min-height:24px;
`
interface DescriptionProps {
  oneLine: boolean;
}
const Description=styled.h5<DescriptionProps>`
  font-weight:normal;
  line-height:20px;
  ${(p)=> p.oneLine&&`
    white-space: nowrap;
    text-overflow: ellipsis;
    overflow:hidden;
  `}
`
interface ImageProps {
  readonly src: string;
}
const Img = styled.div<ImageProps>`
  background-image:url("${(p)=> p.src}");
  background-position:center;
  background-size:cover;
  height:90px;
  width:90px;
  border-radius: 5px;
  position:relative;
`
const Tokens=styled.div`
  display:flex;
  align-items:center;
`
