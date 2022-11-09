import React, { useState } from 'react';
import styled from 'styled-components';
import FadeLeft from '../../animated/fadeLeft';
import { IconButton } from '../../sphinxUI';

export default function GalleryViewer(props) {
  const { gallery, wrap, selectable, big, showAll, style, cover } = props;
  const [selectedImage, setSelectedImage] = useState(0);
  const g = gallery;

  if (!g || !g.length) return <div />;
  //<Square big={big} />

  const showNav = g.length > 1;

  function next(e) {
    e.stopPropagation();
    const nextindex = selectedImage + 1;
    if (g[nextindex]) setSelectedImage(nextindex);
    else setSelectedImage(0);
  }

  function prev(e) {
    e.stopPropagation();
    const previndex = selectedImage - 1;
    if (g[previndex]) setSelectedImage(previndex);
    else setSelectedImage(g.length - 1);
  }

  return (
    <>
      <Gallery
        style={{ width: big || wrap ? '100%' : 'fit-content', ...style }}
        ref={props.innerRef}
      >
        {showAll ? (
          <div style={{ textAlign: 'center' }}>
            {g.map((ga, i) => {
              return <BigImg big={big} src={ga} cover={cover} key={i} />;
            })}
          </div>
        ) : (
          <>
            <Img big={big} src={g[selectedImage]} cover={cover} />
            {showNav && (
              <L>
                <Circ>
                  <IconButton iconStyle={{ color: '#000' }} icon={'chevron_left'} onClick={prev} />
                </Circ>
              </L>
            )}

            {showNav && (
              <R>
                <Circ>
                  <IconButton icon={'chevron_right'} iconStyle={{ color: '#000' }} onClick={next} />
                </Circ>
              </R>
            )}

            <Label>
              {selectedImage + 1} / {g.length}
            </Label>
          </>
        )}

        {/* <FadeLeft isMounted={selectedImage}>
                <BigEnv onClick={() => setSelectedImage('')}>
                    <BigImg src={selectedImage} />
                </BigEnv>
            </FadeLeft> */}
      </Gallery>
    </>
  );
}

const Gallery = styled.div`
  display: flex;
  flex-wrap: wrap;
  position: relative;
`;

const Label = styled.div`
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  border-radius: 20px;
  position: absolute;
  right: 5px;
  bottom: 10px;
  background: rgba(0, 0, 0, 0.25);
  border-radius: 21px;
  font-size: 12px;
  line-height: 18px;
  /* or 152% */
  padding: 4px 10px;
  display: flex;
  align-items: center;
  text-align: center;
  letter-spacing: 1px;
`;

const Circ = styled.div`
  display: flex;
  align-items: center;
  justify-content: center;
  width: 30px;
  height: 30px;
  background: #ffffff66;
  border-radius: 20px;
  cursor: pointer;
`;

const R = styled.div`
  position: absolute;
  right: 5px;
  top: 0px;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
`;

const L = styled.div`
  position: absolute;
  left: 5px;
  top: 0px;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
`;

interface ImageProps {
  readonly src?: string;
  big?: boolean;
  cover?: boolean;
}
const Img = styled.div<ImageProps>`
  background-image: url('${(p) => p.src}');
  background-position: center;
  background-repeat: no-repeat;
  background-size: ${(p) => (p.cover ? 'cover' : 'contain')};
  height: ${(p) => (p.big ? '372px' : '132px')};
  width: ${(p) => (p.big ? '100%' : '132px')};
  position: relative;
`;

const BigImg = styled.img<ImageProps>`
  background-image: url('${(p) => p.src}');
  // background-size: contain;
  background-position: center;
  background-repeat: no-repeat;
  background-size: ${(p) => (p.cover ? 'cover' : 'contain')};
  max-width: 100%;
  height: auto;
  margin-bottom: 5px;
`;
