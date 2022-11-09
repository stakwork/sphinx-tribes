import React, { useState } from 'react';
import styled from 'styled-components';
import { EuiGlobalToastList } from '@elastic/eui';
import MaterialIcon from '@material/react-material-icon';
import { Button } from '../../sphinxUI';
export default function BotBar(props: any) {
  const { value, simple } = props;
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

  function copyToClipboard(str) {
    const el = document.createElement('textarea');
    el.value = str;
    document.body.appendChild(el);
    el.select();
    document.execCommand('copy');
    document.body.removeChild(el);
    addToast();
  }

  return (
    <>
      <Button
        text={`/bot install ${value}`}
        wideButton={true}
        icon={'content_copy'}
        iconStyle={{ color: '#fff', fontSize: 20 }}
        iconSize={20}
        color={'primary'}
        width={'100%'}
        height={50}
        style={{ width: '100%', padding: '0 5px' }}
        onClick={() => copyToClipboard(`/bot install ${value}`)}
      />

      <EuiGlobalToastList toasts={toasts} dismissToast={removeToast} toastLifeTimeMs={1000} />
    </>
  );
}

const QRWrap = styled.div`
  font-family: Roboto;
  font-style: normal;
  font-weight: normal;

  font-family: Roboto;
  font-style: normal;
  font-weight: normal;
  font-size: 15px;
  line-height: 18px;
  /* identical to box height */

  letter-spacing: 0.02em;

  /* Main bottom icons */

  // color: #5F6368;
`;
const Row = styled.div`
  display: flex;
  justify-content: space-between;
  height: 48px;
  width: 100%;
  align-items: center;
  font-family: Roboto;
  font-style: normal;
  font-weight: normal;
  font-size: 15px;
  line-height: 48px;
  /* identical to box height, or 320% */

  display: flex;
  align-items: center;

  /* Secondary Text 4 */

  color: #8e969c;
  cursor: pointer;
  border-radius: 5px;
  // border:1px solid #5F636877;
  // padding: 0 12px 0 6px;
  color: #5f6368;
  &:hover {
    // background:#618AFF44;
    // color:#ffffff;
    // border:1px solid #618AFF44;
  }
`;

const Copy = styled.div`
  font-family: Roboto;
  font-style: normal;
  font-weight: 500;
  font-size: 11px;
  line-height: 13px;
  display: flex;
  align-items: center;
  text-align: right;
  letter-spacing: 0.04em;
  text-transform: uppercase;

  /* Primary blue */

  color: #618aff;
`;
