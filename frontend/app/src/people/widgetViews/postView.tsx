import moment from 'moment';
import React, { useState } from 'react'
import styled from "styled-components";
import { Post } from '../../form/inputs/widgets/interfaces';
import GalleryViewer from '../utils/galleryViewer';
import ReactMarkdown from 'react-markdown'
import { Title, Paragraph, Date, Link } from '../../sphinxUI';

export default function PostView(props: Post) {
    const { title, content, created, gallery } = props
    const isLong = content && (content.length > 100)
    const [expand, setExpand] = useState(false)

    const noGallery = !gallery || !gallery.length

    return <Wrap style={{ maxHeight: expand ? '' : 472 }}>
        <Pad>
            <Title>{title} </Title>
            <Date>{created && moment.unix(created).format('LLL')} </Date>
            <Paragraph style={{
                maxHeight: (noGallery && !expand) ? 340 : !expand ? 80 : '',
                minHeight: noGallery ? 80 : '', overflow: 'hidden'
            }}>
                <ReactMarkdown>{content}</ReactMarkdown> </Paragraph>

            {isLong &&
                <Link
                    style={{ textAlign: 'right', width: '100%' }}
                    onClick={(e) => {
                        e.stopPropagation()
                        setExpand(!expand)
                    }}>
                    {!expand ? 'Read more' : 'Show less'}
                </Link>
            }

        </Pad>
        <GalleryViewer
            showAll={false}
            big={true}
            wrap={false}
            selectable={true}
            gallery={gallery}
            style={{ maxHeight: 291, overflow: 'hidden' }} />

    </Wrap>

}

const Wrap = styled.div`
display: flex;
flex-direction:column;
width:100%;
min-height:472px;
min-width:100%;
font-style: normal;
font-weight: 500;
font-size: 24px;
line-height: 20px;
color: #3C3F41;
justify-content:space-between;

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

// const Link = styled.div`
// font-size: 15px;
// line-height: 20px;
// width:100%;
// text-align:right;
// /* Primary blue */
// color: #618AFF;
// cursor:pointer;
// `;


const M = styled.div`
font-size: 16px;
line-height: 20px;
/* or 167% */


/* Main bottom icons */
text-overflow: ellipsis;
// white-space: nowrap;
color: #5F6368;
margin-bottom:10px;
overflow:hidden;
`;

const Pad = styled.div`
display:flex;
flex-direction:column;
padding:20px;
`;


