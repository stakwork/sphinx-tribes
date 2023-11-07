import React, { useEffect, useState } from 'react';
import styled from 'styled-components';
import lighningDecoder from 'light-bolt11-decoder';
import { InvoiceState } from 'people/widgetViews/organization/interface';
import QR from 'people/utils/QR';
import QrBar from 'people/utils/QrBar';
import { calculateTimeLeft } from '../../../../helpers';
import { QrWrap } from './style';

const InvoiceWrap = styled.div`
  display: flex;
  flex-direction: column;
`;

const CountDownTimerWrap = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 0.3rem;
  margin-bottom: 0.94rem;
`;

const CountDownTextWrapper = styled.div`
  display: flex;
  align-items: center;
`;

const CountDownText = styled.p`
  color: #3c3f41;
  font-family: Barlow;
  font-size: 0.9375rem;
  font-style: normal;
  font-weight: 500;
  line-height: normal;
  margin-bottom: 0;
  margin-left: 0.27rem;
`;

const CountDownIconContainer = styled.div`
  display: flex;
  align-items: center;
  justify-content: center;
  width: 1.5rem;
  height: 1.5rem;
`;

const CountDownTimer = styled.p`
  color: #3c3f41;
  font-family: Barlow;
  font-size: 0.9375rem;
  font-style: normal;
  font-weight: 700;
  line-height: 1.1875rem;
  margin-bottom: 0;
`;

export default function Invoice(props: {
  startDate: Date;
  invoiceStatus: boolean;
  lnInvoice: string;
  invoiceTime: number;
  setInvoiceState?: (state: InvoiceState) => void;
  qrSize?: number;
}) {
  const decoded = lighningDecoder.decode(props.lnInvoice);
  const expiry = decoded.sections[8].value;
  const timeCreated = decoded.sections[4].value;
  const endTime = new Date((timeCreated + expiry) * 1000);
  const [timeLeft, setTimeLeft] = useState(calculateTimeLeft(endTime, 'hours'));

  useEffect(() => {
    const invoiceTimeout = setTimeout(() => {
      setTimeLeft(calculateTimeLeft(endTime, 'hours'));
    }, 1000);

    if (props.invoiceStatus) {
      clearTimeout(invoiceTimeout);
    }

    if (timeLeft.hours === 0 && timeLeft.minutes === 0 && timeLeft.seconds === 0) {
      if (props.setInvoiceState) {
        props.setInvoiceState('EXPIRED');
      }
    }

    return () => {
      if (invoiceTimeout) clearTimeout(invoiceTimeout);
    };
  }, [timeLeft, props.invoiceStatus, props]);

  return (
    <>
      {timeLeft.seconds >= 0 || timeLeft.minutes >= 0 || timeLeft.hours! >= 0 ? (
        <InvoiceWrap>
          <CountDownTimerWrap>
            <CountDownTextWrapper>
              <CountDownIconContainer>
                <img src="/static/count_down.svg" alt="count down" />
              </CountDownIconContainer>
              <CountDownText>Invoice Expires in</CountDownText>
            </CountDownTextWrapper>
            <CountDownTimer>
              {timeLeft.hours}:{timeLeft.minutes}:{timeLeft.seconds}
            </CountDownTimer>
          </CountDownTimerWrap>

          <QrWrap>
            <QR size={props.qrSize || 200} value={props.lnInvoice} />
            <QrBar value={props.lnInvoice} simple style={{ marginTop: '0.94rem' }} />
          </QrWrap>
        </InvoiceWrap>
      ) : null}
    </>
  );
}
