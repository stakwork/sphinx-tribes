import React from 'react'
import styled from "styled-components";

export default function Blog(props: any) {

    return <Wrap>
        <div>{props.label}</div>
        <div>{props.description}</div>
        <div>{props.header}</div>
        <div>{props.title}</div>
    </Wrap>

}

const Wrap = styled.div`
color: #fff;
display: flex;
`;

