import React from 'react'
import styled from "styled-components";
import { Offer } from '../../form/inputs/widgets/interfaces';
import { formatPrice } from '../../helpers';
import { useIsMobile } from '../../hooks';
import GalleryViewer from '../utils/galleryViewer';
import { Divider, Title } from '../../sphinxUI';
import NameTag from '../utils/nameTag';

export default function OfferView(props: any) {
    const { gallery, title, description, price, person, created } = props
    const isMobile = useIsMobile()

    if (isMobile) {
        return <Wrap>

            <GalleryViewer gallery={gallery} selectable={false} wrap={false} big={false} showAll={false} />

            <Body>
                <T>{title || 'No title'}</T>
                <D>{description || 'No description'}</D>
                <P>{formatPrice(price)} <B>SAT</B></P>
            </Body>

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
                <NameTag {...person} created={created} widget={'offer'} />
                <DT>{title}</DT>
                <DD style={{ maxHeight: gallery ? 40 : '' }}>{description}</DD>
            </Pad>
            <Divider style={{ margin: 0 }} />

            <Pad style={{ padding: 20, }}>
                <P style={{ fontSize: 17 }}>{formatPrice(price)} <B>SAT</B></P>
            </Pad>
        </div>
    </DWrap >

}
const Wrap = styled.div`
display: flex;
justify-content:flex-start;

`;

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
margin-bottom:10px;
overflow:hidden;

font-family: Roboto;
font-style: normal;
font-weight: normal;
font-size: 13px;
line-height: 20px;
/* or 154% */

/* Main bottom icons */

color: #5F6368;

`;
const DT = styled(Title)`
margin-bottom:9px;
max-height:52px;
overflow:hidden;

/* Text 2 */

color: #3C3F41;
`;

const T = styled.div`
font-weight:bold;
overflow:hidden;
line-height: 20px;
text-overflow: ellipsis;
display: -webkit-box;
-webkit-line-clamp: 2;
-webkit-box-orient: vertical;
`;
const B = styled.span`
font-weight:300;
`;
const P = styled.div`
font-weight:500;
`;
const D = styled.div`
color:#5F6368;
overflow:hidden;
line-height:18px;
text-overflow: ellipsis;
display: -webkit-box;
-webkit-line-clamp: 2;
-webkit-box-orient: vertical;
`;


const Body = styled.div`
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
height:132px;
`;

const Pad = styled.div`
display:flex;
flex-direction:column;

`;



