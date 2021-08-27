import React from 'react'
import styled from 'styled-components'
import { EuiButton } from '@elastic/eui';
import MaterialIcon from '@material/react-material-icon';
import FadeLeft from '../animated/fadeLeft';
import { IconButton } from '.';

export default function Modal(props: any) {
    const { visible, fill, overlayClick, drift, dismountCallback, children, close, style, hideOverlay } = props

    const fillStyle = fill ? {
        height: '100%',
        width: '100%',
        borderRadius: 0
    } : {}
    return <FadeLeft
        withOverlay={!hideOverlay}
        drift={100}
        direction='up'
        overlayClick={overlayClick}
        dismountCallback={dismountCallback}
        isMounted={visible ? true : false}
        style={{
            position: 'absolute', top: 0, left: 0,
            zIndex: 1000000, width: '100%', height: '100%',
            display: 'flex', alignItems: 'center', justifyContent: 'center',
            ...style
        }}>
        <Env style={{ ...fillStyle }}>
            {close && <X>
                <IconButton
                    onClick={close}
                    size={20}
                    icon='close'
                />
            </X>}
            {children}
        </Env>
    </FadeLeft>
}

const X = styled.div`
position:absolute;
top:5px;
right:0px;
cursor:pointer;
`

const Env = styled.div`
width: 312px;
min-height: 254px;
display: flex;
align-items: center;
justify-content: center;
border-radius: 16px;
background:#ffffff;
position:relative;
`

