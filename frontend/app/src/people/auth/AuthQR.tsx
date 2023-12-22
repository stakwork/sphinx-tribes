import React, { useState, useEffect } from 'react';
import { EuiLoadingSpinner } from '@elastic/eui';
import styled from 'styled-components';
import { observer } from 'mobx-react-lite';
import { AuthProps } from 'people/interfaces';
import { formatRelayPerson } from 'helpers';
import api from '../../api';
import { useStores } from '../../store';
import type { MeInfo } from '../../store/ui';
import { getHost } from '../../config/host';
import { QR } from '../../components/common';

const ConfirmWrap = styled.div`
  display: flex;
  align-items: center;
  justify-content: center;
  flex-direction: column;
  width: 100%;
  height: 203px;
  position: relative;
`;

//TODO: mv to utils
const host = getHost();
function makeQR(challenge: string, ts: string) {
  return `sphinx.chat://?action=auth&host=${host}&challenge=${challenge}&ts=${ts}`;
}

let interval;

function AuthQR(props: AuthProps) {
  const { ui, main } = useStores();
  const [challenge, setChallenge] = useState('');
  const [ts, setTS] = useState('');

  const qrString = makeQR(challenge, ts);

  async function startPolling(challenge: string) {
    let i = 0;
    interval = setInterval(async () => {
      try {
        const me: MeInfo = await api.get(`poll/${challenge}`);
        if (me && me?.pubkey) {
          await ui.setMeInfo(me);
          const person = formatRelayPerson(me);
          await main.saveProfile(person);

          await main.getSelf(me);
          setChallenge('');
          if (props.onSuccess) props.onSuccess();
          if (interval) clearInterval(interval);
        }
        i++;
        if (i > 100) {
          if (interval) clearInterval(interval);
        }
      } catch (e) {
        console.log(e, 'Error');
      }
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
    return function cleanup() {
      if (interval) clearInterval(interval);
    };
  }, []);

  return (
    <ConfirmWrap style={{ ...props.style }}>
      {challenge ? (
        <QR size={203} style={{ width: 203 }} value={qrString} />
      ) : (
        <EuiLoadingSpinner size="xl" />
      )}
    </ConfirmWrap>
  );
}

export default observer(AuthQR);
