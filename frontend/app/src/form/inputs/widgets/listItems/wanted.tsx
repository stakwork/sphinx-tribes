import React from 'react'
import styled from "styled-components";
import * as I from '../interfaces';


export default function Wanted(props: I.Wanted) {

    return <Wrap>
        <div>{props.title}</div>
        <div>{props.description}</div>
        <div>{props.priceMin}</div>
        <div>{props.priceMax}</div>
        <div>{props.url}</div>
    </Wrap>

}

const Wrap = styled.div`
color: #fff;
display: flex;
`;

