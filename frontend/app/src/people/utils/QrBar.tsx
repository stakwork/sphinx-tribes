import React, { useState } from 'react';
import styled from 'styled-components';
import MaterialIcon from '@material/react-material-icon';
import { EuiGlobalToastList } from '@elastic/eui';
import { QRBarProps } from 'people/interfaces';

const QRWrap = styled.div`
  font-family: Roboto;
  font-style: normal;
  font-weight: normal;
  font-size: 13px;
  line-height: 15px;
  letter-spacing: 0.02em;

  /* Main bottom icons */

  color: #5f6368;
`;
const Row = styled.div`
  display: flex;
  justify-content: space-between;
  align-items: center;
  cursor: pointer;
`;

const Value = styled.p`
  color: #5f6368;
  text-overflow: ellipsis;
  font-family: Barlow;
  font-size: 0.8125rem;
  font-style: normal;
  font-weight: 400;
  line-height: normal;
  letter-spacing: 0.01625rem;
  margin-bottom: 0;
`;

const Copy = styled.p`
  color: #5d8fdd;
  text-align: right;
  font-family: Barlow;
  font-size: 0.75rem;
  font-style: normal;
  font-weight: 500;
  line-height: normal;
  letter-spacing: 0.03rem;
  text-transform: uppercase;
  margin-bottom: 0;
`;

export default function QrBar(props: QRBarProps) {
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

  const formatRequest = (request: string) => `${request.substring(0, 23)}...`;

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

  return (
    <Row style={props.style} onClick={() => copyToClipboard(value)}>
      <QRWrap
        style={{
          display: 'flex',
          alignItems: 'center'
        }}
      >
        {!simple && (
          <MaterialIcon
            icon={'qr_code_2'}
            style={{ fontSize: 20, color: '#B0B7BC', marginRight: 10 }}
          />
        )}
        <Value>{formatRequest(value)}</Value>
      </QRWrap>

      <Copy
        style={{
          display: 'flex',
          fontSize: 11,
          alignItems: 'center',
          color: '#618AFF',
          cursor: 'pointer',
          letterSpacing: '0.3px'
        }}
      >
        COPY
      </Copy>

      <EuiGlobalToastList toasts={toasts} dismissToast={removeToast} toastLifeTimeMs={1000} />
    </Row>
  );
}
