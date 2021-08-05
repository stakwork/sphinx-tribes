import React from 'react'
import styled from "styled-components";


export default function TwitterView(props) {

    return <Wrap>
        <div>{props.title}</div>
        <div>{props.createdAt}</div>
    </Wrap>

}

const Wrap = styled.div`
color: #fff;
display: flex;
`;

