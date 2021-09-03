import React from 'react'
import styled from 'styled-components';
import { QRCode } from "react-qr-svg";
// import { QRCode } from 'react-qrcode-logo';
export default function QR(props) {



    return <QRCode
        bgColor="#FFFFFF"
        fgColor="#000000"
        // ecLevel={'Q'}
        level={'Q'}
        // size={203}
        width={props.size}
        height={props.size}
        // logoImage={'static/sphinx.png'}
        // logoWidth={40}
        // logoHeight={40}
        value={props.value}
    />

}
