import React, { useEffect, useState } from 'react';
import QR from 'people/utils/QR';
import QrBar from 'people/utils/QrBar';
import { calculateTimeLeft } from '../../../../helpers';
import { CountDownText, CountDownTimer, CountDownTimerWrap, InvoiceWrap, QrWrap } from './style';

export default function Invoice(props: {
  startDate: Date;
  invoiceStatus: boolean;
  lnInvoice: string;
  invoiceTime: number;
}) {
  const [timeLimit] = useState(props.startDate);
  const [timeLeft, setTimeLeft] = useState(calculateTimeLeft(timeLimit, 'minutes'));

  useEffect(() => {
    const invoiceTimeout = setTimeout(() => {
      setTimeLeft(calculateTimeLeft(timeLimit, 'minutes'));
    }, 1000);

    if (props.invoiceStatus) {
      clearTimeout(invoiceTimeout);
    }
  }, [timeLeft, props.invoiceStatus, timeLimit]);

  return (
    <div style={{ marginTop: '30px' }}>
      {timeLeft.seconds >= 0 || timeLeft.minutes >= 0 ? (
        <InvoiceWrap>
          <CountDownTimerWrap>
            <CountDownText>Invoice expires in {props.invoiceTime} minute</CountDownText>
            <CountDownTimer>
              {timeLeft.minutes}:{timeLeft.seconds}
            </CountDownTimer>
          </CountDownTimerWrap>

          <QrWrap>
            <QR size={220} value={props.lnInvoice} />
            <QrBar value={props.lnInvoice} simple style={{ marginTop: 11 }} />
          </QrWrap>
        </InvoiceWrap>
      ) : null}
    </div>
  );
}
