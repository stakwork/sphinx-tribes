import React, { useState } from 'react'
import styled from "styled-components";
import FadeLeft from '../../animated/fadeLeft';

export default function GalleryViewer(gallery) {
    const [selectedImage, setSelectedImage] = useState('')

    const showImages = gallery && gallery.length
    if (!showImages) return <div />

    return <Gallery>

        {gallery && gallery.map((g, i) => {
            return <Img key={i} src={g} onClick={() => setSelectedImage(g)} />
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
width:100%;
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
                border-radius: 5px;
                position: relative;
                border:1px solid #ffffff21;
                `;