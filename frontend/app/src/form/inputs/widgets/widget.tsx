import React from 'react'
import styled from "styled-components";

export default function Widget(props: any) {

    return <Wrap onClick={() => props.setSelected(props)}>
        {props.icon && <Icon source={props.icon} />}
        <div>{props.label}</div>
    </Wrap>

}

const Wrap = styled.div`
    color: #fff;
    height: 100px;
    width: 200px;
    border: 1px solid #fff;
    border-radius:5px;
    margin-bottom:10px;
    display: flex;
    flex-direction: column;
    align-content: center;
    justify-content: space-evenly;
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

