import React from 'react';
import styled from 'styled-components';
import moment from 'moment';
import { BlogPost } from '../../components/form/inputs/widgets/interfaces';

const Wrap = styled.div`
  display: flex;
  flex-direction: column;
  width: 100%;
  min-width: 100%;
  font-weight: 500;
  font-size: 24px;
  line-height: 20px;

  color: #3c3f41;
`;

const T = styled.div`
  border-radius: 5px;
  font-weight: bold;
  font-size: 25px;
`;

const Time = styled.div`
  border-radius: 5px;
  font-weight: bold;
  font-size: 25px;
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
  background-image: url('${(p: any) => p.src}');
  background-position: center;
  background-size: cover;
  height: 80px;
  width: 80px;
  border-radius: 5px;
  position: relative;
`;

export default function BlogView(props: BlogPost) {
  const { title, markdown, gallery, created } = props;

  const showImages = gallery && gallery.length;
  return (
    <Wrap>
      <T>{title || 'No title'} </T>
      <Time>{moment(created).format('l') || 'No title'} </Time>
      <M>{markdown || 'No markdown'} </M>

      {showImages && (
        <Gallery>{gallery && gallery.map((g: any, i: number) => <Img key={i} src={g} />)}</Gallery>
      )}
    </Wrap>
  );
}
