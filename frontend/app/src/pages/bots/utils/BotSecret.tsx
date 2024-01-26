import React, { useState } from 'react';
import styled from 'styled-components';
import { EuiGlobalToastList } from '@elastic/eui';
import { observer } from 'mobx-react-lite';
import { BotSecretProps } from '../interfaces';
import { Button } from '../../../components/common';
import { useStores } from '../../../store';

const Head = styled.div`
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  width: 100%;
`;

const Name = styled.div`
  width: 100%;
  font-style: normal;
  font-weight: 500;
  /* or 73% */

  text-align: center;

  /* Text 2 */

  color: #3c3f41;
  margin: 20px 0;

  font-size: 24px;
  line-height: 22px;
`;

const RowWrap = styled.div`
  display: flex;
  justify-content: space-between;
  min-height: 48px;
  align-items: center;
  font-family: Roboto;
  font-style: normal;
  font-weight: normal;
  font-size: 15px;
  line-height: 48px;
  width: 100%;
  color: #8e969c;
`;

function BotSecret(props: BotSecretProps) {
  const { ui } = useStores();
  const { meInfo } = ui || {};
  const { id, secret, name, full } = props;
  const [toasts, setToasts]: any = useState([]);

  function addToast() {
    setToasts([
      {
        id: '1',
        title: 'Copied!'
      }
    ]);
  }

  function removeToast() {
    setToasts([]);
  }

  function copyToClipboard(str: string) {
    const el = document.createElement('textarea');
    el.value = str;
    document.body.appendChild(el);
    el.select();
    document.execCommand('copy');
    document.body.removeChild(el);
    addToast();
  }

  function makeToken() {
    const url = `${meInfo?.url}/action`;
    return `${btoa(id)}.${btoa(secret)}.${btoa(url)}`;
  }

  const value = makeToken();

  return (
    <>
      <div
        style={{
          padding: '15px 20px',
          borderRadius: 6,
          border: !full ? '1px dashed #618AFF' : '',
          background: !full ? '#618aff0a' : ''
        }}
      >
        {full && (
          <Head>
            <RowWrap>
              <Name>{name} was created!</Name>
            </RowWrap>
          </Head>
        )}
        <div
          style={{
            fontSize: 15,
            textAlign: 'center',
            width: '100%',
            marginBottom: 20,
            color: '#3C3F41'
          }}
        >
          Use this secret to connect with Sphinx.
        </div>

        <Button
          text={value}
          wideButton={true}
          icon={'content_copy'}
          iconStyle={{ color: '#fff', fontSize: 20 }}
          iconSize={20}
          color={'widget'}
          width={'100%'}
          height={50}
          style={{ width: '100%', padding: '0 5px', fontSize: 12 }}
          onClick={() => copyToClipboard(value)}
        />
      </div>

      <EuiGlobalToastList toasts={toasts} dismissToast={removeToast} toastLifeTimeMs={1000} />
    </>
  );
}

export default observer(BotSecret);
