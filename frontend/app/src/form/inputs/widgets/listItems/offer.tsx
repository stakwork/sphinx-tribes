import React from 'react'
import styled from "styled-components";
import * as I from '../interfaces';


export default function Offer(props: I.Offer) {

    return <Wrap>
        <div>{props.price}</div>
        <div>{props.description}</div>
        <div>{props.img}</div>
        <div>{props.title}</div>
    </Wrap>

}

const Wrap = styled.div`
color: #fff;
display: flex;
`;

