import React, {useState, useEffect} from 'react'
import { QRCode } from 'react-qr-svg';
import {
  EuiLoadingSpinner,
} from '@elastic/eui';
import styled from 'styled-components'
import api from '../api'
import { useStores } from '../store';
import type {Tokens} from '../store/ui'

const host = window.location.host==='localhost:3001'?'localhost:5002':window.location.host
function makeQR(challenge:string) {
  return `sphinx.chat://?action=tokens&host=${host}&challenge=${challenge}`
}

export default function ConfirmMe(){
  const {ui} = useStores()
  const [challenge, setChallenge] = useState('')

  async function startPolling(challenge:string){
    let ok = true
    let i = 0
    while(ok) {
      await sleep(3000)
      try {
        const ts:Tokens = await api.get(`poll/${challenge}`)
        console.log(ts)
        if(ts && ts.pubkey) {
          ui.setTokens(ts)
          ok = false
          break;
        }
        i ++
        if(i>100) ok = false
      } catch(e) {}
    }
  }
  async function getChallenge(){
    const res = await api.get('ask')
    if(res.challenge) {
      setChallenge(res.challenge)
      startPolling(res.challenge)
    }
  }
  useEffect(()=>{
    getChallenge()
  }, [])
  // 1. get "challenge" 
  // 2. display
  // 3. poll for signed timestamp
  // 4. when its gotten, go store in localstorage
  const qrString = makeQR(challenge)
  return <ConfirmWrap>
    {!challenge && <EuiLoadingSpinner size="xl" style={{marginTop:60}} />}
    {challenge && <InnerWrap>
      <QRCode
        bgColor="#FFFFFF"
        fgColor="#000000"
        level="Q"
        style={{width:209}}
        value={qrString}
      />
      <LinkWrap>
        <a href={qrString} className="btn join-btn">
          <img style={{width:13,height:13,marginRight:8}} src="/static/launch-24px.svg" alt="" />
          Open Sphinx
        </a>
      </LinkWrap>
    </InnerWrap>}
  </ConfirmWrap>
}

const ConfirmWrap = styled.div`
  display:flex;
  align-items:center;
  justify-content:center;
  flex-direction:column;
  width:100%;
  min-height:250px;
`
const InnerWrap = styled.div`
  display:flex;
  align-items:center;
  justify-content:center;
  flex-direction:column;
  width:100%;
`
const LinkWrap = styled.div`
  width:100%;
  text-align:center;
  margin:20px 0;
  & a {
    width:115px;
    position:relative;
    margin-left:25px;
  }
`

async function sleep(ms:number) {
	return new Promise(resolve => setTimeout(resolve, ms))
}