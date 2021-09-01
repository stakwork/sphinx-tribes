import React from 'react'
import styled from "styled-components";
import { Wanted } from '../../form/inputs/widgets/interfaces';
import { formatPrice } from '../../helpers';
import { useIsMobile } from '../../hooks';
import GalleryViewer from '../utils/galleryViewer';
import { Divider } from '../../sphinxUI';

export default function WantedView(props: Wanted) {
    const { title, description, priceMin, priceMax, url, gallery } = props
    const isMobile = useIsMobile()

    if (isMobile) {
        return <Wrap>

            <GalleryViewer gallery={gallery} selectable={false} wrap={false} big={false} showAll={false} />

            <Body>
                <T>{title || 'No title'}</T>
                <D>{description || 'No description'}</D>
                <P>{formatPrice(priceMin)} <B>sat</B> - {formatPrice(priceMax)} <B>sat</B></P>
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
            style={{ maxHeight: 291, overflow: 'hidden' }} />

        <Pad style={{ padding: 20 }}>
            <DT>{title || 'No title'}</DT>
            <DD>{description || 'No description'}</DD>
        </Pad>
        <Divider style={{ margin: 0 }} />
        <Pad style={{ padding: 20, }}>
            <P style={{ fontSize: 17 }}>{formatPrice(priceMin)} <B>sat</B> - {formatPrice(priceMax)} <B>sat</B></P>
        </Pad>
    </DWrap>

}

const DWrap = styled.div`
display: flex;
flex:1;
height:100%;
min-height:100%;
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
padding:10px;
`;


const DD = styled.div`
font-style: normal;
font-weight: normal;
font-size: 12px;
line-height: 25px;
margin-bottom:10px;


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