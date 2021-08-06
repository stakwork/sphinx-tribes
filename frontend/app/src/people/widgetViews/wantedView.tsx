import React from 'react'
import styled from "styled-components";
import { Wanted } from '../../form/inputs/widgets/interfaces';
import { formatPrice } from '../../helpers';

export default function WantedView(props: Wanted) {
    const { title, description, priceMin, priceMax, url } = props

    return < Wrap >
        <P>{formatPrice(priceMin)} ~ {formatPrice(priceMax)}</P>
        <T>{title || 'No title'}</T>
        <Body>
            <D>{description || 'No description'}</D>
            <U>{url || 'No link'}</U>
        </Body>

    </Wrap >

}

interface ImageProps {
    readonly src: string;
}
const Img = styled.div<ImageProps>`
                background-image: url("${(p) => p.src}");
                background-position: center;
                background-size: cover;
                height: 80px;
                width: 80px;
                border-radius: 5px;
                position: relative;
                border:1px solid #ffffff21;
                `;

const T = styled.div`
font-weight:bold;
`;
const P = styled.div`
padding: 10px;
max-width: 400;
background: #ffffff21;
border-radius: 5px;
text-align:center;
margin-bottom:5px;
`;
const D = styled.div`
`;
const U = styled.div`
color:#1BA9F5
`;
const Wrap = styled.div`
color: #fff;
display: flex;
flex-direction:column;
`;

const Body = styled.div`
padding: 0 5px 5px 5px;
font-size:14px;
`;

const Gallery = styled.div`
display:flex;
flex-wrap:wrap;
margin-top:10px;
`;

const None = styled.div`
color:#ffffff71;
`;

