import moment from 'moment';
import React from 'react'
import styled from "styled-components";
import { Post } from '../../../form/inputs/widgets/interfaces';
import GalleryViewer from '../../utils/galleryViewer';
import ReactMarkdown from 'react-markdown'

export default function PostSummary(props: Post) {
    const { title, content, created, gallery } = props

    return <Pad>
        <T>{title || 'No title'} </T>
        <Time>{created && moment.unix(created).format('LLL')} </Time>
        <M><ReactMarkdown>{content}</ReactMarkdown></M>

        {/* readmore */}

        <GalleryViewer gallery={gallery} showAll={true} wrap={false} selectable={false} big={true} />
    </Pad>

}

const Pad = styled.div`
width:100%;
padding:20px;
`;

const Time = styled.div`
font-style: normal;
font-weight: normal;
font-size: 12px;
line-height: 25px;
margin-bottom:10px;


/* Main bottom icons */

color: #8E969C;
`;
const T = styled.div`
font-style: normal;
font-weight: normal;
font-size: 24px;
line-height: 25px;
color:#292C33 !important;
margin-bottom:10px;

color: #5F6368;
`;

const M = styled.div`
font-size: 15px;
line-height: 25px;
/* or 167% */


/* Main bottom icons */

color: #5F6368;
margin-bottom:20px;
`;

