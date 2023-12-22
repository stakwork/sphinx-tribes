import React, { useState } from 'react';
import { EuiGlobalToastList } from '@elastic/eui';
import { Button } from '../../../components/common';

export default function BotBar(props: { value: string }) {
  const { value } = props;
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
