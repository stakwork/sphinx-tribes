import moment from 'moment';
import React from 'react'
import styled from "styled-components";
import { Post } from '../../form/inputs/widgets/interfaces';


export default function PostView(props: Post) {
    const { title, content, created, gallery } = props

    return <Wrap>
        <T>{title || 'No title'} </T>
        <Time>{created && moment.unix(created).format('LLL')} </Time>
        <M>{content || 'No content'} </M>

        {/* readmore */}

        {<Gallery>
            {gallery && gallery.map((g, i) => {
                return <Img key={i} src={g} />
            })}
        </Gallery>}
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

const M = styled.div`
font-size: 15px;
line-height: 25px;
/* or 167% */


/* Main bottom icons */

color: #5F6368;
`;

const Gallery = styled.div`
display:flex;
margin-top:10px;
`;

interface ImageProps {
    readonly src: string;
}
const Img = styled.div<ImageProps>`
                background-image: url("${(p) => p.src}");
                background-position: center;
                background-size: cover;
                min-height:200px;
                width: 100%;
                border-radius: 5px;
                position: relative;
                border:1px solid #ffffff21;
                `;
