import React from 'react'
import styled from "styled-components";
import { SupportMe } from '../../form/inputs/widgets/interfaces';


export default function SupportMeView(props: SupportMe) {

    return <Wrap>
        <div>{props.title}</div>
        <div>{props.createdAt}</div>
    </Wrap>

}

const Wrap = styled.div`
color: #fff;
display: flex;
`;

