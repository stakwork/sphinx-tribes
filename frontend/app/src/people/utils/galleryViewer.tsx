import React, { useState } from 'react'
import styled from "styled-components";
import FadeLeft from '../../animated/fadeLeft';
import { IconButton } from '../../sphinxUI';

export default function GalleryViewer({ gallery, wrap, selectable, big }) {
    const [selectedImage, setSelectedImage] = useState(0)

    let g = gallery

    if (!g || !g.length) return <div />
    //<Square big={big} />

    const showNav = (g.length > 1) && big

    function next(e) {
        e.stopPropagation()
        let nextindex = selectedImage + 1
        if (g[nextindex]) setSelectedImage(nextindex)
        else setSelectedImage(0)
    }

    function prev(e) {
        e.stopPropagation()
        let previndex = selectedImage - 1
        if (g[previndex]) setSelectedImage(previndex)
        else setSelectedImage(g.length - 1)
    }

    return <>
        <Gallery style={{ width: (big || wrap) ? '100%' : 'fit-content', minHeight: big && 370 }}>

            <Img big={big} src={g[selectedImage]} />
            {showNav &&
                <L>
                    <Circ>
                        <IconButton
                            iconStyle={{ color: '#000' }}
                            icon={'chevron_left'}
                            onClick={prev}
                        />
                    </Circ>
                </L>
            }

            {showNav &&
                <R>
                    <Circ>
                        <IconButton
                            icon={'chevron_right'}
                            iconStyle={{ color: '#000' }}
                            onClick={next}
                        />
                    </Circ>
                </R>
            }

            <Label>
                {selectedImage + 1} / {g.length}
            </Label>


            {/* <FadeLeft isMounted={selectedImage}>
                <BigEnv onClick={() => setSelectedImage('')}>
                    <BigImg src={selectedImage} />
                </BigEnv>
            </FadeLeft> */}

        </Gallery>
    </>
}

const Gallery = styled.div`
display:flex;
flex-wrap:wrap;
position:relative;
`;

const Label = styled.div`
display:flex;
align-items:center;
justify-content:center;
color:#fff;
border-radius:20px;
position:absolute;
right:5px;
bottom:10px;
background: rgba(0, 0, 0, 0.25);
border-radius: 21px;
font-size: 12px;
line-height: 18px;
/* or 152% */
padding:4px 10px;
display: flex;
align-items: center;
text-align: center;
letter-spacing: 1px;
`;

const Circ = styled.div`
display:flex;
align-items:center;
justify-content:center;
width:30px;
height:30px;
background:#ffffff99;
border-radius:20px;
cursor:pointer;
`

const R = styled.div`
position:absolute;
right:5px;
top:0px;
height:100%;
display:flex;
align-items:center;
justify-content:center;
`;

const L = styled.div`
position:absolute;
left:5px;
top:0px;
height:100%;
display:flex;
align-items:center;
justify-content:center;
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