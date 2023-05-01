import React, { useEffect, useState } from 'react';
import { CountDownText, CountDownTimer, CountDownTimerWrap } from './style';
import { useStores } from '../../../../store';
import QR from 'people/utils/QR';

export default function Invoice(props: { startDate: Date, count: number }) {
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

    return (
        <div style={{ marginTop: "30px" }}>
            {timeLeft.seconds ?
                <>
                    <CountDownTimerWrap>
                        <CountDownText>Invoice expires in a minute</CountDownText>
                        <CountDownTimer>{timeLeft.minutes}:{timeLeft.seconds}</CountDownTimer>
                    </CountDownTimerWrap>

                    <QR size={220} value={main.lnInvoice} />
                </>
                : null}
        </div>
    )
}