import React from 'react'
import styled from "styled-components";
import { BlogPost } from '../../form/inputs/widgets/interfaces';


export default function BlogView(props: BlogPost) {
    const { title, markdown, gallery } = props

    const showImages = gallery && gallery.length
    return <Wrap>
        <T>{title || 'No title'} </T>
        <M>{markdown || 'No markdown'} </M>

        {showImages && <Gallery>
            {gallery && gallery.map((g, i) => {
                return <Img key={i} src={g} />
            })}
        </Gallery>}
    </Wrap>

}

const Wrap = styled.div`
color: #fff;
display: flex;
flex-direction:column;
`;

const T = styled.div`
padding: 10px;
background: #ffffff21;
border-radius: 5px;
text-align:center;
margin-bottom:5px;
font-weight:bold;
`;

const M = styled.div`
border-radius: 5px;
margin:15px 0 10px 0;
`;
const Gallery = styled.div`
display:flex;
flex-wrap:wrap;
margin-top:10px;
`;

interface ImageProps {
    readonly src: string;
}
const Img = styled.div<ImageProps>`
                background-image: url("${(p) => p.src}");
                background-position: center;
                background-size: cover;
                height: 80px;
                width: 80px;
                border-radius: 5px;
                position: relative;
                border:1px solid #ffffff21;
                `;