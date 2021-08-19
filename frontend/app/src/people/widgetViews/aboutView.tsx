import React from 'react'
import styled from "styled-components";
import { Divider } from '../../sphinxUI';

export default function AboutView(props: any) {
    const { price_to_meet, description, extras } = props

    console.log('AboutView', price_to_meet)

    const { twitter } = extras || {}
    const { handle } = twitter || {}

    return <Wrap>
        <Row>
            <div>Price to Join:</div>
            <div>{price_to_meet}</div>
        </Row>

        <Divider />

        <D>{description || 'No description'} </D>

        <T>Follow Me</T>
        {handle && <div>@{handle}</div>}

        {/* show twitter etc. here */}
    </Wrap>

}



const Wrap = styled.div`
display: flex;
flex-direction:column;
`;
const Row = styled.div`
display:flex;
justify-content:space-between;
margin-bottom:20px;

`;
const T = styled.div`
font-weight:bold;
margin-bottom:5px;
`;

const D = styled.div`

margin:15px 0 10px 0;
`;

