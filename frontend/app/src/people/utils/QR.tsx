import React from 'react'
import styled from 'styled-components';
import { QRCode } from "react-qr-svg";

export default function QR(props) {

    return <div style={{ position: 'relative' }}>
        <QRCode
            bgColor="#FFFFFF"
            fgColor="#000000"
            level={'Q'}
            style={{ width: props.size }}
            value={props.value}
        />

        {/* logo env */}
        <div style={{
            position: 'absolute', zIndex: 10,
            height: props.size,
            width: props.size,
            top: 0, left: 0,
            display: 'flex', justifyContent: 'center', alignItems: 'center'
        }}>
            <Img src={'/static/sphinx.png'} />
        </div>
    </div>

}

interface ImageProps {
    readonly src: string;
}
const Img = styled.div<ImageProps>`
            background-image: url("${(p) => p.src}");
            background-position: center;
            background-size: cover;
            height: 55px;
            width: 55px;
            border-radius: 50%;
            `;