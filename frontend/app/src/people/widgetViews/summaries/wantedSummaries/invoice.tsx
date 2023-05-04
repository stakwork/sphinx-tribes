import React, { useEffect, useState } from 'react';
import { CopyInvoiceBtn, CountDownText, CountDownTimer, CountDownTimerWrap, InvoiceWrap } from './style';
import { useStores } from '../../../../store';
import QR from 'people/utils/QR';

export default function Invoice(props: { startDate: Date, count: number, dataStatus: boolean }) {
    const [timeLimit] = useState(props.startDate);

    const calculateTimeLeft = () => {
        const difference = new Date(timeLimit).getTime() - new Date().getTime();

        let timeLeft: any = {};

        if (difference > 0) {
            timeLeft = {
                minutes: Math.floor((difference / 1000 / 60) % 60),
                seconds: Math.floor((difference / 1000) % 60),
            };
        }

        return timeLeft;
    };


    const { main } = useStores();
    const [timeLeft, setTimeLeft] = useState(calculateTimeLeft());

    useEffect(() => {
        const invoiceTimeout = setTimeout(() => {
            setTimeLeft(calculateTimeLeft());
        }, 1000);

        if (props.count > 29) {
            clearTimeout(invoiceTimeout);
        }

    }, [timeLeft, props.count]);

    const copyInvoice = () => {
        navigator.clipboard.writeText(main.lnInvoice)
    }

    return (
        <div style={{ marginTop: "30px" }}>
            {timeLeft.seconds && !props.dataStatus ?
                <InvoiceWrap>
                    <CountDownTimerWrap>
                        <CountDownText>Invoice expires in a minute</CountDownText>
                        <CountDownTimer>{timeLeft.minutes}:{timeLeft.seconds}</CountDownTimer>
                    </CountDownTimerWrap>

                    <QR size={220} value={main.lnInvoice} />

                    <CopyInvoiceBtn onClick={copyInvoice}>Copy invoice</CopyInvoiceBtn>
                </InvoiceWrap>
                : null}
        </div>
    )
}