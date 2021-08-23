import React, { useState } from 'react'
import styled from "styled-components";
import FadeLeft from '../../animated/fadeLeft';

export default function GalleryViewer({ gallery, wrap, selectable, big }) {
    const [selectedImage, setSelectedImage] = useState('')

    let g = gallery

    if (!g || !g.length) return <Square big={big} />

    return <Gallery style={{ width: (big || wrap) ? '100%' : 'fit-content', minHeight: big && 370 }}>
        {g && g.length && g.map((ga, i) => {
            return <Img big={big} key={i} src={ga} onClick={() => {
                if (selectable) setSelectedImage(ga)
            }} />
        })}

        <FadeLeft isMounted={selectedImage}>
            <BigEnv onClick={() => setSelectedImage('')}>
                <BigImg src={selectedImage} />
            </BigEnv>
        </FadeLeft>

    </Gallery>
}

const Gallery = styled.div`
display:flex;
flex-wrap:wrap;
`;

const Square = styled.div<ImageProps>`
    background-position: center;
    background-size: cover;
    height: ${(p) => p.big ? '372px' : '132px'};
    width: ${(p) => p.big ? '100%' : '132px'};
    position: relative;
    background:#5F636822;
`;

const BigEnv = styled.div`
position:absolute;
top:0px;
left:0px;
display:flex;
justify-content:center;
align-items:center;
z-index:20;
`;

interface ImageProps {
    readonly src?: string;
    big?: boolean;
}
const Img = styled.div<ImageProps>`
                background-image: url("${(p) => p.src}");
                background-position: center;
                background-size: cover;
                height: ${(p) => p.big ? '372px' : '132px'};
                width: ${(p) => p.big ? '100%' : '132px'};
                // border-radius: 5px;
                position: relative;
                border:1px solid #ffffff21;
                `;

const BigImg = styled.div<ImageProps>`
                background-image: url("${(p) => p.src}");
                background-position: center;
                background-size: cover;
                height: 60vh;
                min-height:200px;
                max-height:800px;
                width: 60vw;
                min-width:200px;
                max-width:800px;
                // border-radius: 5px;
                position: relative;
                border:1px solid #ffffff21;
                `;