import React from 'react'
import styled from "styled-components";
import { Post } from '../../form/inputs/widgets/interfaces';


export default function PostView(props: Post) {
    const { title, content } = props

    return <Wrap>
        <T>{title || 'No title'} </T>
        <M>{content || 'No content'} </M>
    </Wrap>

}

const Wrap = styled.div`
display: flex;
flex-direction:column;
width:100%;
min-width:100%;
`;

const T = styled.div`
border-radius: 5px;
font-weight:bold;
font-size:25px;
`;

const M = styled.div`
border-radius: 5px;
margin:15px 0 10px 0;
`;