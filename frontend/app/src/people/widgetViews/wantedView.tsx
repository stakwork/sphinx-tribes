import React from 'react'
import styled from "styled-components";
import { Wanted } from '../../form/inputs/widgets/interfaces';
import { formatPrice } from '../../helpers';
import { useIsMobile } from '../../hooks';
import GalleryViewer from '../utils/galleryViewer';
import { Divider } from '../../sphinxUI';

export default function WantedView(props: any) {
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
            style={{ maxHeight: 276, overflow: 'hidden' }} />

        <Pad style={{ padding: 20 }}>
            <DT>{title || 'No title'}</DT>
            <DD style={{ maxHeight: gallery ? 40 : '' }}>{description || 'No description'}</DD>
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
const DT = styled.div`
margin-bottom:9px;
font-family: Roboto;
font-style: normal;
font-weight: 500;
font-size: 15px;
line-height: 20px;
/* or 133% */
max-height:40px;
overflow:hidden;

/* Primary Text 1 */

color: #292C33;
`;