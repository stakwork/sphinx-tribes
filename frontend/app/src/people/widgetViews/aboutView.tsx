import React from 'react'
import styled from "styled-components";
import { Divider } from '../../sphinxUI';

export default function AboutView(props: any) {
    const { price_to_meet, description, extras } = props

    const { twitter } = extras || {}
    let tag = ''
    if (twitter && twitter[0] && twitter[0].value) tag = twitter[0].value
    return <Wrap>
        <Row>
            <div>Price to Join:</div>
            <div>{price_to_meet}</div>
        </Row>

        <Divider />

        <D>{description || 'No description'} </D>

        {tag && <>
            <T>Follow Me</T>
            <I>
                <Icon source={`/static/twitter.png`} />
                <div>@{tag}</div>
            </I>
        </>}
        {/* <I>Facebook</I> */}
        {/* <div></div>
        {handle && <div>@{handle}</div>} */}

        {/* show twitter etc. here */}
    </Wrap>

}



const Wrap = styled.div`
display: flex;
flex-direction:column;
`;
const I = styled.div`
display:flex;
align-items:center;

`;
const Row = styled.div`
display:flex;
justify-content:space-between;
margin-bottom:20px;

`;
const T = styled.div`
font-weight:bold;
margin-top:5px;
margin-bottom:5px;
`;

const D = styled.div`

margin:15px 0 10px 0;
`;

interface IconProps {
    source: string;
}

const Icon = styled.div<IconProps>`
                    background-image: ${p => `url(${p.source})`};
                    width:20px;
                    height:20px;
                    margin-right:5px;
                    background-position: center; /* Center the image */
                    background-repeat: no-repeat; /* Do not repeat the image */
                    background-size: contain; /* Resize the background image to cover the entire container */
                    border-radius:5px;
                    overflow:hidden;
                `;