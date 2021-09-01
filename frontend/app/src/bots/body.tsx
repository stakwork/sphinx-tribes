import React, { useEffect, useState } from 'react'
import styled from 'styled-components';

export default function BotsBody() {

    return <div style={{
        display: 'flex', flexDirection: 'column',
        justifyContent: 'center', alignItems: 'center',
        height: '100%', background: '#f0f1f3',
    }}>
        <Icon src={'static/coming_soon.png'} />

        <>
            <H>COMING SOON</H>
            <C>Stay tuned for something amazing!</C>
        </>

        <div style={{ height: 200 }} />
    </div>


}


interface IconProps {
    src: string;
}

const Icon = styled.div<IconProps>`
                    background-image: ${p => `url(${p.src})`};
                    width:160px;
                    height:160px;
                    margin-right:10px;
                    background-position: center; /* Center the image */
                    background-repeat: no-repeat; /* Do not repeat the image */
                    background-size: contain; /* Resize the background image to cover the entire container */
                    border-radius:5px;
                    overflow:hidden;
                `;

const H = styled.div`
margin-top:10px;
font-size: 30px;
font-style: normal;
font-weight: 700;
line-height: 37px;
letter-spacing: 0.1em;
text-align: center;
            `;

const C = styled.div`
margin-top:10px;
font-family: Roboto;
font-size: 22px;
font-style: normal;
font-weight: 400;
line-height: 26px;
letter-spacing: 0em;
text-align: center;
color: #8E969C;


            `;