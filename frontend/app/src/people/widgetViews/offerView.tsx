import React from 'react'
import styled from "styled-components";
import { Offer } from '../../form/inputs/widgets/interfaces';
import { formatPrice } from '../../helpers';
import GalleryViewer from '../utils/galleryViewer';

export default function OfferView(props: Offer) {
    const { gallery, title, description, price, url } = props

    return <Wrap>

        <T>{title || 'No title'}</T>
        <Body>
            <D>{description || 'No description'}</D>
            <U>{url || 'No link'}</U>
            <GalleryViewer gallery={gallery} />
        </Body>
        <P>{formatPrice(price)} </P>

    </Wrap>

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
                `;

const T = styled.div`
font-weight:bold;
`;
const P = styled.div`
padding: 10px;
max-width: 400;
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
display: flex;
flex-direction:column;
`;

const Body = styled.div`
padding: 0 5px 5px 5px;
font-size:14px;
`;
