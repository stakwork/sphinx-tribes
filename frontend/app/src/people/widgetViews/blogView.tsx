import React from 'react'
import styled from "styled-components";
import { BlogPost } from '../../form/inputs/widgets/interfaces';


export default function BlogView(props: BlogPost) {

    return <Wrap>
        <div>{props.title}</div>
        <div>{props.createdAt}</div>
    </Wrap>

}

const Wrap = styled.div`
color: #fff;
display: flex;
`;

