import React from 'react'
import styled from "styled-components";
import { Divider } from '../../sphinxUI';
import QrBar from '../utils/QrBar'
import ReactMarkdown from 'react-markdown'

export function renderMarkdown(str) {
    return <ReactMarkdown>{str}</ReactMarkdown>
}

export default function AboutView(props: any) {
    const { price_to_meet, description, extras, twitter_confirmed, owner_pubkey } = props
    const { twitter } = extras || {}
    let tag = ''
    if (twitter && twitter[0] && twitter[0].value) tag = twitter[0].value

    return <Wrap>
        <Row>
            <div>Price to Connect:</div>
            <div style={{ fontWeight: 'bold', color: '#000' }}>{price_to_meet}</div>
        </Row>

        <Divider />

        <QrBar value={owner_pubkey} />

        <Divider />

        <D>{renderMarkdown(description)}</D>

        {tag && <>
            <T>For Normies</T>
            <I>
                <Icon source={`/static/twitter2.png`} />
                <Tag>@{tag}</Tag>
                {twitter_confirmed ?
                    <Badge>VERIFIED</Badge> :
                    <Badge style={{ background: '#b0b7bc' }}>PENDING</Badge>
                }

            </I>
        </>}
        {/* <I>Facebook</I> */}
        {/* <div></div>
        {handle && <div>@{handle}</div>} */}

        {/* show twitter etc. here */}


    </Wrap>

}
const Badge = styled.div`
display:flex;
justify-content:center;
align-items:center;
margin-left:10px;
height:20px;
color:#ffffff;
background: #1DA1F2;
border-radius: 32px;
font-weight: bold;
font-size: 8px;
line-height: 9px;
padding 0 10px;
`;
const QRWrap = styled.div`
font-family: Roboto;
font-style: normal;
font-weight: normal;
font-size: 13px;
line-height: 15px;
letter-spacing: 0.02em;

/* Main bottom icons */

color: #5F6368;
`;
const Wrap = styled.div`
display: flex;
flex-direction:column;
width:100%;
`;
const I = styled.div`
display:flex;
align-items:center;

`;

const Tag = styled.div`
font-family: Roboto;
font-style: normal;
font-weight: normal;
font-size: 15px;
line-height: 26px;
/* or 173% */

display: flex;
align-items: center;

/* Main bottom icons */

color: #5F6368;
`

const Row = styled.div`
display:flex;
justify-content:space-between;
height:48px;
align-items:center;
font-family: Roboto;
font-style: normal;
font-weight: normal;
font-size: 15px;
line-height: 48px;
/* identical to box height, or 320% */

display: flex;
align-items: center;

/* Secondary Text 4 */

color: #8E969C;

`;
const T = styled.div`
font-family: Roboto;
font-style: normal;
font-weight: bold;
font-size: 10px;
line-height: 26px;
/* or 260% */

letter-spacing: 0.3px;
text-transform: uppercase;

/* Text 2 */

color: #3C3F41;
margin-top:5px;
margin-bottom:5px;
`;

const D = styled.div`

margin:15px 0 10px 0;
font-family: Roboto;
font-style: normal;
font-weight: normal;
font-size: 15px;
line-height: 20px;
/* or 133% */


/* Main bottom icons */

color: #5F6368;

`;

interface IconProps {
    source: string;
}

const Icon = styled.div<IconProps>`
                    background-image: ${p => `url(${p.source})`};
                    width:16px;
                    height:13px;
                    margin-right:8px;
                    background-position: center; /* Center the image */
                    background-repeat: no-repeat; /* Do not repeat the image */
                    background-size: contain; /* Resize the background image to cover the entire container */
                    border-radius:5px;
                    overflow:hidden;
                `;