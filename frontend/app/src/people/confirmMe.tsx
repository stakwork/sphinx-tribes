import React, {useState} from 'react'
import { QRCode } from 'react-qr-svg';
import {
  EuiLoadingSpinner,
} from '@elastic/eui';
import styled from 'styled-components'

export default function ConfirmMe({setMe}:{setMe:Function}){
  const [challenge, setChallenge] = useState('')
  // 1. get "challenge" 
  // 2. display
  // 3. poll for signed timestamp
  // 4. when its gotten, go store in localstorage
  return <ConfirmWrap>
    {!challenge && <EuiLoadingSpinner size="xl" />}
    {challenge && <QRCode
      bgColor="#FFFFFF"
      fgColor="#000000"
      level="Q"
      style={{width:209}}
      value={challenge}
    />}
  </ConfirmWrap>
}

const ConfirmWrap = styled.div`
  padding-top:80px;
  display:flex;
  align-items:center;
  justify-content:center;
  width:100%;
`