import React from 'react';
import styled from 'styled-components';
import { SupportMe } from '../../form/inputs/widgets/interfaces';

export default function SupportMeView(props: SupportMe) {
  const { gallery, description, url } = props;

  const showImages = gallery && gallery.length;
  return (
    <Wrap>
      <T>Support Me</T>
      <M>{description || 'No description'} </M>
      <U>{url || 'No link'}</U>

      {showImages && (
        <Gallery>
          {gallery &&
            gallery.map((g, i) => {
              return <Img key={i} src={g} />;
            })}
        </Gallery>
      )}
    </Wrap>
  );
}

const Wrap = styled.div`
  display: flex;
  flex-direction: column;
`;
const U = styled.div`
  color: #1ba9f5;
`;
const T = styled.div`
  padding: 10px;
  background: #ffffff21;
  border-radius: 5px;
  text-align: center;
  margin-bottom: 5px;
  // font-weight:bold;
`;

const M = styled.div`
  border-radius: 5px;
  margin: 15px 0 10px 0;
`;
const Gallery = styled.div`
  display: flex;
  flex-wrap: wrap;
  margin-top: 10px;
`;

interface ImageProps {
  readonly src: string;
}
const Img = styled.div<ImageProps>`
  background-image: url('${(p) => p.src}');
  background-position: center;
  background-size: cover;
  height: 80px;
  width: 80px;
  border-radius: 5px;
  position: relative;
`;
