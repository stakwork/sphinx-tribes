import React from 'react';
import { Button, QR } from '../../components/common';

export interface TorSaveQRProps {
  url: string;
  goBack: Function;
}

export default function TorSaveQR(props: TorSaveQRProps) {
  const { url, goBack } = props;

  return (
    <div
      style={{
        display: 'flex',
        flexDirection: 'column',
        justifyContent: 'center',
        alignItems: 'center',
        padding: '10px 20px',
        width: '100%'
      }}
    >
      <div style={{ height: 40 }} />
      <div style={{ textAlign: 'center', margin: '0 0 20px' }}>Scan to complete this request.</div>
      <QR size={220} value={url} />

      <Button
        text={'Open Sphinx App'}
        height={60}
        style={{ marginTop: 30 }}
        width={'100%'}
        color={'primary'}
        onClick={() => {
          const el = document.createElement('a');
          el.href = url;
          el.click();
        }}
      />

      <Button
        text={'Cancel'}
        height={60}
        style={{ marginTop: 20 }}
        width={'100%'}
        color={'action'}
        onClick={() => {
          goBack();
        }}
      />
    </div>
  );
}
