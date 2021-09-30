import moment from 'moment';
import React from 'react'
import styled from "styled-components";
import { Post } from '../../../form/inputs/widgets/interfaces';
import GalleryViewer from '../../utils/galleryViewer';
import ReactMarkdown from 'react-markdown'
import NameTag from '../../utils/nameTag';
import MaterialIcon from '@material/react-material-icon';
import FavoriteButton from '../../utils/favoriteButton';

export function renderMarkdown(str) {
    return <ReactMarkdown>{str}</ReactMarkdown>
}

export default function PostSummary(props: any) {
    const { title, content, created, gallery, person } = props

    const heart = <FavoriteButton />

    return <Pad>
        <div style={{ display: 'flex', width: '100%', justifyContent: 'space-between' }}>
            <NameTag {...person} created={created} widget={'post'} />
            {heart}
        </div>

        <T>{title || 'No title'} </T>
        {/* <Time>{created && moment.unix(created).format('LLL')} </Time> */}
        <M>{renderMarkdown(content)}</M>

        {/* readmore */}

        <GalleryViewer gallery={gallery} showAll={true} wrap={false} selectable={false} big={true} />
    </Pad>

}

const Pad = styled.div`
width:100%;
max-width:602px;
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

