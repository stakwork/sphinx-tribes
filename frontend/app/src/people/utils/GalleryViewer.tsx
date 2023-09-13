import React, { useState } from 'react';
import styled from 'styled-components';
import { IconButton } from '../../components/common';

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
  background-image: url('${(p: any) => p.src}');
  background-position: center;
  background-repeat: no-repeat;
  background-size: ${(p: any) => (p.cover ? 'cover' : 'contain')};
  height: ${(p: any) => (p.big ? '372px' : '132px')};
  width: ${(p: any) => (p.big ? '100%' : '132px')};
  position: relative;
`;

const BigImg = styled.img<ImageProps>`
  background-image: url('${(p: any) => p.src}');
  // background-size: contain;
  background-position: center;
  background-repeat: no-repeat;
  background-size: ${(p: any) => (p.cover ? 'cover' : 'contain')};
  max-width: 100%;
  height: auto;
  margin-bottom: 5px;
`;

export default function GalleryViewer(props: any) {
  const { gallery, wrap, big, showAll, style, cover } = props;
  const [selectedImage, setSelectedImage] = useState(0);
  const g = gallery;

  if (!g || !g.length) return <div />;
  //<Square big={big} />

  const showNav = g.length > 1;

  function next(e: any) {
    e.stopPropagation();
    const nextindex = selectedImage + 1;
    if (g[nextindex]) setSelectedImage(nextindex);
    else setSelectedImage(0);
  }

  function prev(e: any) {
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
            {g.map((ga: any, i: number) => (
              <BigImg big={big} src={ga} cover={cover} key={i} />
            ))}
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
      </Gallery>
    </>
  );
}
