import React, { useEffect, useState } from 'react';
import { CountDownText, CountDownTimer, CountDownTimerWrap, InvoiceWrap, QrWrap } from './style';
import { useStores } from '../../../../store';
import QR from 'people/utils/QR';
import { calculateTimeLeft } from '../../../../helpers';
import QrBar from 'people/utils/QrBar';
import { invoicePollTarget } from 'config';

export default function Invoice(props: {
  startDate: Date;
  count: number;
  dataStatus?: boolean;
  pollMinutes: number;
}) {
  const [timeLimit] = useState(props.startDate);

  const { main } = useStores();
  const [timeLeft, setTimeLeft] = useState(calculateTimeLeft(timeLimit, 'minutes'));

  useEffect(() => {
    const invoiceTimeout = setTimeout(() => {
      setTimeLeft(calculateTimeLeft(timeLimit, 'minutes'));
    }, 1000);

    if (props.count > invoicePollTarget * props.pollMinutes) {
      clearTimeout(invoiceTimeout);
    }
  }, [timeLeft, props.count]);

  return (
    <div style={{ marginTop: '30px' }}>
      {timeLeft.seconds >= 0 || (timeLeft.minutes >= 0) ? (
        <InvoiceWrap>
          <CountDownTimerWrap>
            <CountDownText>Invoice expires in a minute</CountDownText>
            <CountDownTimer>
              {timeLeft.minutes}:{timeLeft.seconds}
            </CountDownTimer>
          </CountDownTimerWrap>

          <QrWrap>
            <QR size={220} value={main.lnInvoice} />
            <QrBar value={main.lnInvoice} simple style={{ marginTop: 11 }} />
          </QrWrap>
        </InvoiceWrap>
      ) : null}
    </div>
  );
}
