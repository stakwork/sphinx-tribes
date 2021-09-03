import React, { useState } from 'react'
import styled from "styled-components";
import { Button, Divider } from '../../sphinxUI';
import MaterialIcon from '@material/react-material-icon';
import {
    EuiGlobalToastList,
} from '@elastic/eui';


export default function QrBar(props: any) {
    const { value, simple } = props
    const [toasts, setToasts]: any = useState([]);

    function addToast() {
        setToasts([{
            id: '1',
            title: 'Copied!'
        }]);
    };

    function removeToast() {
        setToasts([]);
    };


    function copyToClipboard(str) {
        const el = document.createElement('textarea');
        el.value = str;
        document.body.appendChild(el);
        el.select();
        document.execCommand('copy');
        document.body.removeChild(el);
        addToast()
    };

    return <Row style={props.style}>
        <QRWrap style={{
            display: 'flex', alignItems: 'center', width: '70%',
            overflow: 'hidden',
            whiteSpace: 'nowrap',
            textOverflow: 'ellipsis',
        }}>
            {!simple && <MaterialIcon
                icon={'qr_code_2'}
                style={{ fontSize: 20, color: '#B0B7BC', marginRight: 10 }} />}

            <div style={{
                overflow: 'hidden',
                whiteSpace: 'nowrap',
                textOverflow: 'ellipsis'
            }}>
                {value}
            </div>
        </QRWrap>

        <Button
            text='Copy'
            color={'link2'}
            // width={72}
            // height={32}
            // style={{ minWidth: 72 }}
            onClick={() => copyToClipboard(value)}
        />

        <EuiGlobalToastList
            toasts={toasts}
            dismissToast={removeToast}
            toastLifeTimeMs={1000}
        />
    </Row>

}
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