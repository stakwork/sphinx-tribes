import React, { useState, useEffect } from 'react';
import { EuiLoadingSpinner } from '@elastic/eui';
import styled from 'styled-components';
import api from '../api';
import { useStores } from '../store';
import type { MeInfo } from '../store/ui';
import { getHost } from '../host';

const host = getHost();
function makeQR(challenge: string, ts: string) {
  return `sphinx.chat://?action=auth&host=${host}&challenge=${challenge}&ts=${ts}`;
}

let interval;

export default function ConfirmMe(props: any) {
  const { ui, main } = useStores();
  const [challenge, setChallenge] = useState('');
  const [ts, setTS] = useState('');

  // const isMobile = useIsMobile()

  const qrString = makeQR(challenge, ts);

  useEffect(() => {
    getChallenge();
  }, []);

  useEffect(() => {
    if (challenge && ts) {
      const el = document.createElement('a');
      el.href = qrString;
      el.click();
    }
  }, [challenge, ts]);

  async function startPolling(challenge: string) {
    let i = 0;
    interval = setInterval(async () => {
      try {
        const me: MeInfo = await api.get(`poll/${challenge}`);
        // console.log(me);
        if (me && me.pubkey) {
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
      } catch (e) { }
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
    <ConfirmWrap>
      <InnerWrap>
        <div style={{ marginBottom: 50 }}>Opening Sphinx...</div>
        <EuiLoadingSpinner size="xl" />
      </InnerWrap>
    </ConfirmWrap>
  );
}

const ConfirmWrap = styled.div`
  display: flex;
  align-items: center;
  justify-content: center;
  flex-direction: column;
  width: 100%;
  min-height: 250px;
`;
const InnerWrap = styled.div`
  display: flex;
  align-items: center;
  justify-content: center;
  flex-direction: column;
  width: 100%;
`;
const LinkWrap = styled.div`
  width: 100%;
  text-align: center;
  margin: 20px 0;
  & a {
    width: 115px;
    position: relative;
    margin-left: 25px;
  }
`;
const P = styled.p`
  margin-top: 10px;
`;
const QrWrap = styled.div`
  padding: 8px;
  background: white;
  display: flex;
  align-items: center;
  justify-content: center;
`;

async function sleep(ms: number) {
  return new Promise((resolve) => setTimeout(resolve, ms));
}
