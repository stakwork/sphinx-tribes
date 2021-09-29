import React from 'react'
import styled from "styled-components";
import { formatPrice } from '../../helpers';
import { useIsMobile } from '../../hooks';
import GalleryViewer from '../utils/galleryViewer';
import { Divider, Title } from '../../sphinxUI';
import NameTag from '../utils/nameTag';
import MaterialIcon from '@material/react-material-icon';

export default function WantedView(props: any) {
    const { title, description, priceMin, priceMax, price, url, gallery, person, created, issue, repo, type } = props
    const isMobile = useIsMobile()

    function getMobileView() {
        return <Wrap>
            <GalleryViewer gallery={gallery} selectable={false} wrap={false} big={false} showAll={false} />
            <Body>
                <T>{title}</T>
                <D>{description}</D>
                <P>{formatPrice(priceMin)} <B>sat</B> - {formatPrice(priceMax)} <B>sat</B></P>
            </Body>
        </Wrap>
    }

    function getDesktopView() {
        if (type === 'coding_task') {
            // this is a code task
            return <DWrap>
                <Pad style={{ padding: 20, height: 410 }}>
                    <div style={{ display: 'flex', width: '100%', justifyContent: 'space-between' }}>
                        <NameTag {...person} created={created} widget={'wanted'} />
                        <MaterialIcon icon={'code'} />
                    </div>
                    <DT>{title}</DT>
                    <DD>{description}</DD>
                    <DD>{repo + '/' + issue}</DD>
                </Pad>
                <Divider style={{ margin: 0 }} />
                <Pad style={{ padding: 20, }}>
                    <P style={{ fontSize: 17 }}>{formatPrice(price)} <B>SAT</B></P>
                </Pad>
            </DWrap>
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
                <NameTag {...person} created={created} widget={'wanted'} />
                <DT>{title}</DT>
                <DD style={{ maxHeight: gallery ? 40 : '' }}>{description}</DD>
            </Pad>
            <Divider style={{ margin: 0 }} />
            <Pad style={{ padding: 20, }}>
                <P style={{ fontSize: 17 }}>{formatPrice(priceMin)} <B>sat</B> - {formatPrice(priceMax)} <B>sat</B></P>
            </Pad>
        </DWrap>
    }

    if (isMobile) {
        return getMobileView()
    }

    return getDesktopView()
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
const DT = styled(Title)`
margin-bottom:9px;
max-height:52px;
overflow:hidden;

/* Primary Text 1 */

color: #292C33;
`;