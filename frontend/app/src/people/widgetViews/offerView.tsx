import React from 'react'
import styled from "styled-components";
import { Offer } from '../../form/inputs/widgets/interfaces';
import { formatPrice } from '../../helpers';
import GalleryViewer from '../utils/galleryViewer';

export default function OfferView(props: Offer) {
    const { gallery, title, description, price } = props

    return <Wrap>

        <GalleryViewer gallery={gallery} selectable={false} wrap={false} big={false} showAll={false} />
        <Pad>
            <Body>
                <T>{title || 'No title'}</T>
                <D>{description || 'No description'}</D>
                <P>{formatPrice(price)} <B>sat</B></P>
            </Body>
        </Pad>
    </Wrap>

}
const Wrap = styled.div`
display: flex;
justify-content:flex-start;

`;
const T = styled.div`
font-weight:bold;
`;
const B = styled.span`
font-weight:300;
`;
const P = styled.div`
font-weight:500;
`;
const D = styled.div`
color:#5F6368;
white-space: nowrap;
height:26px;
text-overflow: ellipsis;
overflow:hidden;
`;


const Body = styled.div`
font-size:14px;
font-size: 15px;
line-height: 20px;
/* or 133% */
padding:10px;
display: flex;
flex-direction:column;
justify-content: space-around;

/* Primary Text 1 */

color: #292C33;
overflow:hidden;
`;

const Pad = styled.div`
display:flex;
flex-direction:column;

`;
