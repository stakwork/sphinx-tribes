import React, { useState, useEffect } from 'react';
import { EuiLoadingSpinner } from '@elastic/eui';
import styled from 'styled-components';
import api from '../../api';
import { useStores } from '../../store';
import type { MeInfo } from '../../store/ui';
import { getHost } from '../../host';
import QR from '../utils/QR';

const host = getHost();
function makeQR(challenge: string, ts: string) {
  return `sphinx.chat://?action=auth&host=${host}&challenge=${challenge}&ts=${ts}`;
}

let interval;

export default function AuthQR(props: any) {
  const { ui, main } = useStores();
  const [challenge, setChallenge] = useState('');
  const [ts, setTS] = useState('');

  const qrString = makeQR(challenge, ts);

  useEffect(() => {
    getChallenge();
    return function cleanup() {
      if (interval) clearInterval(interval);
    };
  }, []);

  async function startPolling(challenge: string) {
    let i = 0;
    interval = setInterval(async () => {
      try {
        const me: MeInfo = await api.get(`poll/${challenge}`);
        console.log(me);
        if (me && me?.pubkey) {
          ui.setMeInfo(me);
          await main.getSelf(me);
          setChallenge('');
          if (props.onSuccess) props.onSuccess();
          if (interval) clearInterval(interval);
        }
        i++;
        if (i > 100) {
          if (interval) clearInterval(interval);
        }
      } catch (e) {}
    }, 3000);
  }
  async function getChallenge() {
    const res = await api.get('ask');
    if (res.challenge) {
      setChallenge(res.challenge);
      startPolling(res.challenge);
    }
    if (res.ts) {
      setTS(res.ts);
    }
  }
  return (
    <ConfirmWrap style={{ ...props.style }}>
      {/* <InnerWrap>
                <QrWrap> */}
      {challenge ? (
        <QR size={203} style={{ width: 203 }} value={qrString} />
      ) : (
        <EuiLoadingSpinner size="xl" />
      )}
      {/* </QrWrap>
            </InnerWrap> */}
    </ConfirmWrap>
  );
}

const ConfirmWrap = styled.div`
  display: flex;
  align-items: center;
  justify-content: center;
  flex-direction: column;
  width: 100%;
  height: 203px;
  position: relative;
`;
const InnerWrap = styled.div`
  display: flex;
  align-items: center;
  justify-content: center;
  flex-direction: column;
  width: 100%;
  height: 100%;
`;

const QrWrap = styled.div`
  padding: 8px;
  background: white;
  display: flex;
  align-items: center;
  justify-content: center;
  height: 100%;
`;

async function sleep(ms: number) {
  return new Promise((resolve) => setTimeout(resolve, ms));
}
