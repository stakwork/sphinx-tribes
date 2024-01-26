import React, { useState, useEffect } from 'react';
import { EuiLoadingSpinner } from '@elastic/eui';
import styled from 'styled-components';
import { AuthProps } from '../../people/interfaces';
import { formatRelayPerson } from '../../helpers';
import api from '../../api';
import { useStores } from '../../store';
import type { MeInfo } from '../../store/ui';
import { getHost } from '../../config/host';

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

const host = getHost();

function makeQR(challenge: string, ts: string) {
  return `sphinx.chat://?action=auth&host=${host}&challenge=${challenge}&ts=${ts}`;
}

let interval;

export default function SphinxAppLoginDeeplink(props: AuthProps) {
  const { ui, main } = useStores();
  const [challenge, setChallenge] = useState('');
  const [ts, setTS] = useState('');

  const qrString = makeQR(challenge, ts);

  async function startPolling(challenge: string) {
    let i = 0;
    interval = setInterval(async () => {
      try {
        const me: MeInfo = await api.get(`poll/${challenge}`);
        if (me && me.pubkey) {
          await ui.setMeInfo(me);
          const person = formatRelayPerson(me);
          await main.saveProfile(person);

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

  return (
    <ConfirmWrap>
      <InnerWrap>
        <div style={{ marginBottom: 50 }}>Opening Sphinx...</div>
        <EuiLoadingSpinner size="xl" />
      </InnerWrap>
    </ConfirmWrap>
  );
}
