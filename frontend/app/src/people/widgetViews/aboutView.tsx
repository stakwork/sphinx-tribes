import React, { useState } from 'react'
import styled from "styled-components";
import { Button, Divider } from '../../sphinxUI';
import MaterialIcon from '@material/react-material-icon';
import { useStores } from '../../store';

export default function AboutView(props: any) {
    const { price_to_meet, description, extras, twitter_confirmed, owner_contact_key, canEdit } = props
    const { twitter } = extras || {}
    let tag = ''
    if (twitter && twitter[0] && twitter[0].value) tag = twitter[0].value

    const { main } = useStores()

    const [showSettings, setShowSettings] = useState(false)

    function copyToClipboard(str) {
        const el = document.createElement('textarea');
        el.value = str;
        document.body.appendChild(el);
        el.select();
        document.execCommand('copy');
        document.body.removeChild(el);
    };

    return <Wrap>
        <Row>
            <div>Price to Connect:</div>
            <div style={{ fontWeight: 'bold', color: '#000' }}>{price_to_meet}</div>
        </Row>

        <Divider />

        <Row>
            <QRWrap style={{
                display: 'flex', alignItems: 'center', width: '70%',
                overflow: 'hidden',
                whiteSpace: 'nowrap',
                textOverflow: 'ellipsis',
            }}>
                <MaterialIcon
                    icon={'qr_code_2'}
                    style={{ fontSize: 20, color: '#B0B7BC', marginRight: 10 }} />
                <div style={{
                    overflow: 'hidden',
                    whiteSpace: 'nowrap',
                    textOverflow: 'ellipsis'
                }}>
                    {owner_contact_key}
                </div>
            </QRWrap>

            <Button
                text='Copy'
                color='widget'
                width={72}
                height={32}
                style={{ minWidth: 72 }}
                onClick={() => copyToClipboard(owner_contact_key)}
            />
        </Row>

        <Divider />

        <D>{description || 'No description'} </D>

        {tag && <>
            <T>Follow Me</T>
            <I>
                <Icon source={`/static/twitter2.png`} />
                <div>@{tag}</div>
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

        {canEdit && <div style={{ cursor: 'pointer', marginTop: 60, fontSize: 12, marginBottom: 20 }} onClick={() => setShowSettings(!showSettings)}>Show Settings</div>}

        {showSettings &&
            <Button
                text={'Delete my account'}
                color={'danger'}
                onClick={() => main.deleteProfile()}
            />
        }
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
font-weight:bold;
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
                    width:20px;
                    height:20px;
                    margin-right:10px;
                    background-position: center; /* Center the image */
                    background-repeat: no-repeat; /* Do not repeat the image */
                    background-size: contain; /* Resize the background image to cover the entire container */
                    border-radius:5px;
                    overflow:hidden;
                `;