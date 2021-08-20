import React from 'react'
import styled from "styled-components";
import { BlogPost } from '../interfaces';


export default function Blog(props: BlogPost) {

    return <Wrap>
        <div>{props.title}</div>
        <div>{props.createdAt}</div>
    </Wrap>

}

const Wrap = styled.div`
display: flex;
`;

