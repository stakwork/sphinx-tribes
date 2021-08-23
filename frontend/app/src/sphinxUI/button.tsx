import React from 'react'
import styled from 'styled-components'
import { EuiButton } from '@elastic/eui';
import MaterialIcon from '@material/react-material-icon';

export default function Button(props: any) {

    const colors = {
        primary: {
            background: '#618AFF',
            color: '#fff'
        },
        white: {
            background: '#fff',
            color: '#5F6368',
            border: '1px solid #DDE1E5'
        },
        link: {
            background: '#fff',
            color: '#618AFF',
            border: '1px solid #A3C1FF'
        },
        widget: {
            background: '#DDE1E5',
            color: '#3C3F41',
        },
    }

    return <B
        style={{ ...colors[props.color], padding: props.icon && '0 0 0 15px', height: props.height, width: props.width, ...props.style }}
        disabled={props.disabled}
        onClick={props.onClick}
    >
        <>
            {props.icon &&
                <div style={{
                    display: 'flex', alignItems: 'center',
                    position: 'absolute', top: 0, left: 3, height: '100%'
                }}>
                    <MaterialIcon
                        icon={props.icon}
                        style={{ fontSize: props.iconSize ? props.iconSize : 30 }} />
                </div>
            }
            <div style={{ display: 'flex', alignItems: 'center' }}>
                {props.leadingIcon && <MaterialIcon
                    icon={props.leadingIcon}
                    style={{
                        fontSize: props.iconSize ? props.iconSize : 20,
                        marginRight: 10
                    }} />
                }
                <>{props.text}</>
            </div>
        </>
    </B>
}

const B = styled(EuiButton)`
position:relative;
border-radius: 100px;
height:36px;
font-weight:bold;
border:none;
font-weight: 500;
font-size: 15px;
line-height: 18px;
display: flex;
align-items: center;
text-align: center;
box-shadow:none !important;
text-transform:none !important;
`

