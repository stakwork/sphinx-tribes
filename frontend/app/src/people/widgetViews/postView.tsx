import moment from 'moment';
import React, { useState } from 'react'
import styled from "styled-components";
import { Post } from '../../form/inputs/widgets/interfaces';
import GalleryViewer from '../utils/galleryViewer';


export default function PostView(props: Post) {
    const { title, content, created, gallery } = props
    const isLong = content && (content.length > 200)
    const [expand, setExpand] = useState(false)


    return <Wrap>
        <Pad>
            <T>{title || 'No title'} </T>
            <Time>{created && moment.unix(created).format('LLL')} </Time>
            <M style={{ maxHeight: !expand ? 120 : '' }}>{content || 'No content'} </M>

            {isLong &&
                <Link onClick={(e) => {
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
            style={{ maxHeight: 208, overflow: 'hidden' }} />

    </Wrap>

}

const Wrap = styled.div`
display: flex;
flex-direction:column;
width:100%;
min-width:100%;
font-style: normal;
font-weight: 500;
font-size: 24px;
line-height: 20px;
/* or 83% */


/* Text 2 */

color: #3C3F41;

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

const Link = styled.div`
font-size: 15px;
line-height: 20px;


/* Primary blue */

color: #618AFF;
margin-bottom:10px;
`;


const M = styled.div`
font-size: 16px;
line-height: 20px;
/* or 167% */


/* Main bottom icons */

color: #5F6368;
margin-bottom:10px;
overflow:hidden;
`;

const Pad = styled.div`
display:flex;
flex-direction:column;
padding:10px;
`;


