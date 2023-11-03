import React, { useEffect, useState } from 'react';
import styled from 'styled-components';
import QR from 'people/utils/QR';
import QrBar from 'people/utils/QrBar';
import { calculateTimeLeft } from '../../../../helpers';
import { QrWrap } from './style';

const InvoiceWrap = styled.div`
  display: flex;
  flex-direction: column;
  width: 12.5rem;
  margin-left: 4.75rem;
  margin-right: 4.75rem;
`;

const CountDownTimerWrap = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
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
  margin-left: 0.37rem;
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
  leading-trim: both;
  text-edge: cap;
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
  }, [timeLeft, props.invoiceStatus]);

  return (
    <>
      {timeLeft.seconds >= 0 || timeLeft.minutes >= 0 ? (
        <InvoiceWrap>
          <CountDownTimerWrap>
            <CountDownTextWrapper>
              <CountDownIconContainer>
                <img src="/static/count_down.svg" />
              </CountDownIconContainer>
              <CountDownText>Invoice Expires in</CountDownText>
            </CountDownTextWrapper>
            <CountDownTimer>
              {timeLeft.minutes}:{timeLeft.seconds}
            </CountDownTimer>
          </CountDownTimerWrap>

          <QrWrap>
            <QR size={200} value={props.lnInvoice} />
            <QrBar value={props.lnInvoice} simple style={{ marginTop: '0.94rem' }} />
          </QrWrap>
        </InvoiceWrap>
      ) : null}
    </>
  );
}
