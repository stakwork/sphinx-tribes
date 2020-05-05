import React, {useState} from 'react'
import { QRCode } from 'react-qr-svg';
import styled from 'styled-components'
import {
  EuiCheckableCard,
  EuiSpacer,
  EuiHighlight,
} from '@elastic/eui';
import Tag from './tag'

export default function Tribe({uuid,name,img,tags,description,selected,select}:any){
  const showTags = tags&&tags.length&&tags.length>0?true:false
  return <>
    <EuiCheckableCard
      id={uuid}
      label={name}
      name={name}
      value={uuid}
      checked={selected}
      onChange={() => select(uuid)}>
      <Content onClick={() => select(selected?'':uuid)} style={{
        height:selected?256:100
      }} selected={selected}>
        <Left style={selected?{maxWidth:'calc(100% - 256px)'}:{}}>
          <Row>
            <Img src={img} />
            <Left style={{padding:'22px 22px 20px',maxHeight:100,maxWidth:'calc(100% - 100px)'}}>
              <Row style={selected?{flexDirection:'column',alignItems:'flex-start'}:{}}>
                <Title style={selected?{color:'white',fontSize:32,minHeight:34}:{}}>
                  {name}
                </Title>
                {showTags && <Tokens style={selected?{marginTop:16}:{}}>
                  {tags.map((t:string)=> <Tag type={t} key={t} />)}
                </Tokens>}
              </Row>
              {!selected && <>
                <EuiSpacer size="m" />
                <Description oneLine={true} style={{minHeight:20}}>
                  {description}
                </Description>
              </>}
            </Left>
          </Row>
          {selected && <div style={{padding:'6px 17px'}}>
            <EuiSpacer size="s" />
            <Description oneLine={false}>
              {description}
            </Description>
          </div>}
        </Left>
        {selected && <QRCode
          bgColor={selected?'#FFFFFF':'#666'}
          fgColor="#000000"
          level="Q"
          style={{width:256}}
          value={uuid}
        />}
      </Content>
    </EuiCheckableCard>
    <EuiSpacer size="m" />
  </>
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
    color:#ccc;
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
  overflow: scroll;
  max-height:142px;
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
  height:100px;
  width:100px;
  position:relative;
`
const Tokens=styled.div`
  display:flex;
  align-items:center;
`