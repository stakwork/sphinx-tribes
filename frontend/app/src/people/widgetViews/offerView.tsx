import React from 'react'
import styled from "styled-components";
import { Offer } from '../../form/inputs/widgets/interfaces';
import { formatPrice } from '../../helpers';
import { useIsMobile } from '../../hooks';
import GalleryViewer from '../utils/galleryViewer';
import { Divider } from '../../sphinxUI';

export default function OfferView(props: Offer) {
    const { gallery, title, description, price } = props
    const isMobile = useIsMobile()

    if (isMobile) {
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

    return <DWrap>
        <GalleryViewer
            showAll={false}
            big={true}
            wrap={false}
            selectable={true}
            gallery={gallery}
            style={{ maxHeight: 276, overflow: 'hidden' }} />
        <div>
            <Pad style={{ padding: 20, height: gallery ? '' : 411 }}>
                <DT>{title || 'No title'}</DT>
                <DD style={{ height: gallery ? 50 : '' }}>{description || 'No description'}</DD>
            </Pad>
            <Divider style={{ margin: 0 }} />

            <Pad style={{ padding: 20, }}>
                <P style={{ fontSize: 17 }}>{formatPrice(price)} <B>sat</B></P>
            </Pad>
        </div>
    </DWrap >

}
const Wrap = styled.div`
display: flex;
justify-content:flex-start;`


const DWrap = styled.div`
display: flex;
flex-direction:column;
width:100%;
min-width:100%;
font-style: normal;
font-weight: 500;
font-size: 24px;
line-height: 20px;
color: #3C3F41;
justify-content:space-between;
`;

const DD = styled.div`
font-style: normal;
font-weight: normal;
font-size: 12px;
line-height: 25px;
margin-bottom:10px;
overflow:hidden;

/* Main bottom icons */

color: #8E969C;
`;
const DT = styled.div`
font-style: normal;
font-weight: normal;
font-size: 24px;
line-height: 25px;
color:#292C33 !important;
margin-bottom:10px;

color: #5F6368;
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



