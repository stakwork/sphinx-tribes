import React from 'react'
import styled from "styled-components";
import { Offer } from '../../form/inputs/widgets/interfaces';


export default function OfferView(props: Offer) {

    return <Wrap>
        <div>{props.title}</div>
        <div>{props.createdAt}</div>
    </Wrap>

}

const Wrap = styled.div`
color: #fff;
display: flex;
`;

