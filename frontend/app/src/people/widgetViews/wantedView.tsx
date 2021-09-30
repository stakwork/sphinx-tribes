import React from 'react'
import styled from "styled-components";
import { formatPrice } from '../../helpers';
import { useIsMobile } from '../../hooks';
import GalleryViewer from '../utils/galleryViewer';
import { Divider, Title } from '../../sphinxUI';
import NameTag from '../utils/nameTag';
import MaterialIcon from '@material/react-material-icon';
import { extractGithubIssue } from '../../helpers';

export default function WantedView(props: any) {
    const { title, description, priceMin, priceMax, price, url, gallery, person, created, issue, repo, type } = props
    const isMobile = useIsMobile()

    function renderCodingTask() {
        const { assignee, status } = extractGithubIssue(person, repo, issue)

        if (isMobile) {
            return <Wrap>
                <Body style={{ width: '100%' }}>
                    <div style={{ display: 'flex', width: '100%', justifyContent: 'space-between' }}>
                        <NameTag {...person} created={created} widget={'wanted'} style={{ margin: 0 }} />
                        <MaterialIcon icon={'code'} />
                    </div>
                    <T style={{ marginBottom: 5 }}>{title}</T>
                    <div style={{ width: '100%', display: 'flex', justifyContent: 'space-between' }}>
                        <div style={{ display: 'flex', marginBottom: 0 }}>
                            <Status>{status || 'Open'} -</Status>
                            <Assignee>{assignee ? 'Assigned' : 'Unassigned'}</Assignee>
                        </div>
                        <P>{formatPrice(price)} <B>SAT</B></P>
                    </div>
                </Body>
            </Wrap>
        }

        return <DWrap>
            <Pad style={{ padding: 20, height: 410 }}>
                <div style={{ display: 'flex', width: '100%', justifyContent: 'space-between' }}>
                    <NameTag {...person} created={created} widget={'wanted'} />
                    <MaterialIcon icon={'code'} />
                </div>

                <DT style={{ marginBottom: 4 }}>{title}</DT>
                <div style={{ display: 'flex', marginBottom: 0 }}>
                    <Status>{status || 'Open'} -</Status>
                    <Assignee>{assignee ? 'Assigned' : 'Unassigned'}</Assignee>
                </div>
                <Link >github.com/{repo + '/' + issue}</Link>
                <div style={{ height: 15 }} />
                <DescriptionCodeTask style={{ height: 240 }}>{description}</DescriptionCodeTask>

            </Pad>
            <Divider style={{ margin: 0 }} />
            <Pad style={{ padding: 20, }}>
                <P style={{ fontSize: 17 }}>{formatPrice(price)} <B>SAT</B></P>
            </Pad>
        </DWrap>
    }

    function getMobileView() {
        return <Wrap>
            <GalleryViewer gallery={gallery} selectable={false} wrap={false} big={false} showAll={false} />
            <Body>
                <NameTag {...person} created={created} widget={'wanted'} style={{ margin: 0 }} />
                <T>{title}</T>
                <D>{description}</D>
                <P>{formatPrice(priceMin) || '0'} <B>SAT</B> - {formatPrice(priceMax)} <B>SAT</B></P>
            </Body>
        </Wrap>
    }

    function getDesktopView() {
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
                <P style={{ fontSize: 17 }}>{formatPrice(priceMin) || '0'} <B>SAT</B> - {formatPrice(priceMax)} <B>SAT</B></P>
            </Pad>
        </DWrap>
    }

    if (type === 'coding_task') {
        return renderCodingTask()
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

const Assignee = styled.div`
display: flex;
font-size:12px;
font-weight:300;
`;

const Status = styled.div`
display: flex;
font-size:12px;
margin-right:4px;
font-weight:300;
`;

const Link = styled.div`
color:blue;
overflow-wrap:break-word;
font-size:15px;
font-weight:300;
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

overflow: hidden;
text-overflow: ellipsis;
display: -webkit-box;
-webkit-line-clamp: 2;
-webkit-box-orient: vertical;

`;

const DescriptionCodeTask = styled.div`
margin-bottom:10px;

font-family: Roboto;
font-style: normal;
font-weight: normal;
font-size: 13px;
line-height: 20px;
color: #5F6368;
overflow: hidden;
text-overflow: ellipsis;
display: -webkit-box;
-webkit-line-clamp: 12;
-webkit-box-orient: vertical;
height: 240px;
`
const DT = styled(Title)`
margin-bottom:9px;
max-height:52px;
overflow:hidden;

/* Primary Text 1 */

color: #292C33;
`;