import React from 'react'
import styled from "styled-components";

export default function Widget(props: any) {

    return <Wrap onClick={() => props.setSelected(props)}>
        {/* {props.icon && <Icon source={props.icon} />} */}
        <Content>{props.label}</Content>
    </Wrap>

}

const Content = styled.div`
    width:100%;
    display:flex;
    align-content: center;
    justify-content: center;
`

const Wrap = styled.div`
    color: #fff;
    height: 50px;
    width: 100px;
    border: 1px solid #fff;
    border-radius:5px;
    margin:5px;

    display: flex;
    flex-direction: column;
    align-content: center;
    justify-content: center;
`;

export interface IconProps {
    source: string;
}

const Icon = styled.img<IconProps>`
    background-image: ${p => `url(${p.source})`};
    width:100px;
    height:100px;
    background-position: center; /* Center the image */
    background-repeat: no-repeat; /* Do not repeat the image */
    background-size: contain; /* Resize the background image to cover the entire container */
`;

