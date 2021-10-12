import React from 'react'
import styled from "styled-components";
import { formatPrice, satToUsd } from '../../helpers';
import { useIsMobile } from '../../hooks';
import GalleryViewer from '../utils/galleryViewer';
import { Divider, Title } from '../../sphinxUI';
import NameTag from '../utils/nameTag';
import MaterialIcon from '@material/react-material-icon';
import { extractGithubIssue } from '../../helpers';
import GithubStatusPill from './parts/statusPill';

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
                        <Img src={'/static/github_logo.png'} />
                        {/* <MaterialIcon icon={'code'} /> */}
                    </div>
                    <T style={{ marginBottom: 5 }}>{title}</T>
                    <div style={{ width: '100%', display: 'flex', justifyContent: 'space-between' }}>
                        <GithubStatusPill status={status} assignee={assignee} />
                        <P>{formatPrice(price)} <B>SAT ({satToUsd(price)})</B> </P>
                    </div>
                </Body>
            </Wrap>
        }

        return <DWrap>
            <Pad style={{ padding: 20, height: 410 }}>
                <div style={{ display: 'flex', width: '100%', justifyContent: 'space-between' }}>
                    <NameTag {...person} created={created} widget={'wanted'} />
                    {/* <MaterialIcon icon={'code'} /> */}
                    <Img src={'/static/github_logo.png'} />
                </div>

                <DT>{title}</DT>

                <Link >github.com/{repo + '/issues/' + issue}</Link>
                <GithubStatusPill status={status} assignee={assignee} style={{ marginTop: 10 }} />

                <div style={{ height: 15 }} />
                <DescriptionCodeTask>{description}</DescriptionCodeTask>

            </Pad>
            <Divider style={{ margin: 0 }} />
            <Pad style={{ padding: 20, }}>
                <P style={{ fontSize: 17 }}>{formatPrice(price)} <B>SAT ({satToUsd(price)})</B> </P>
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

    if (type === 'coding_task' || type === 'wanted_coding_task') {
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
font-weight:300;
`;

const Pill = styled.div`
display: flex;
justify-content:center;
align-items:center;
height:20px;
font-size:12px;
font-weight:300;
background:green;
border-radius:30px;

border: 1px solid transparent;

padding: 5px 12px;
font-size: 14px;
font-weight: 500;
line-height: 20px;
white-space: nowrap;
border-radius: 2em;
height:32px;
width:81px;
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
-webkit-line-clamp: 8;
-webkit-box-orient: vertical;
height: 160px;
`
const DT = styled(Title)`
margin-bottom:9px;
max-height:52px;
overflow:hidden;

/* Primary Text 1 */

color: #292C33;
`;

interface ImageProps {
    readonly src?: string;
}
const Img = styled.div<ImageProps>`
                        background-image: url("${(p) => p.src}");
                        background-position: center;
                        background-size: cover;
                        position: relative;
                        width:22px;
                        height:22px;
                        `;