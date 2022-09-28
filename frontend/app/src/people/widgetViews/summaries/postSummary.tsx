import React from 'react';
import styled from 'styled-components';
import GalleryViewer from '../../utils/galleryViewer';
import ReactMarkdown from 'react-markdown';
import NameTag from '../../utils/nameTag';
import FavoriteButton from '../../utils/favoriteButton';

export function renderMarkdown(str) {
  return <ReactMarkdown>{str}</ReactMarkdown>;
}

export default function PostSummary(props: any) {
  const { title, content, created, gallery, person } = props;

  const heart = <FavoriteButton />;

  return (
    <div style={{ padding: '40px 20px', overflow: 'auto' }}>
      <Pad>
        <div style={{ display: 'flex', width: '100%', justifyContent: 'space-between' }}>
          <NameTag {...person} iconSize={24} textSize={13} created={created} widget={'post'} />
          {heart}
        </div>

        <T>{title || 'No title'} </T>
        <M>{renderMarkdown(content)}</M>

        <GalleryViewer
          gallery={gallery}
          showAll={true}
          wrap={false}
          selectable={false}
          big={true}
        />
      </Pad>
    </div>
  );
}

const Pad = styled.div`
  width: 100%;
  max-width: 602px;
  padding: 20px;
`;

const Time = styled.div`
  font-style: normal;
  font-weight: normal;
  font-size: 12px;
  line-height: 25px;
  margin-bottom: 10px;

  /* Main bottom icons */

  color: #8e969c;
`;
const T = styled.div`
  font-style: normal;
  font-weight: normal;
  font-size: 24px;
  line-height: 25px;
  color: #292c33 !important;
  margin-bottom: 10px;

  color: #5f6368;
`;

const M = styled.div`
  font-size: 15px;
  line-height: 25px;
  /* or 167% */

  /* Main bottom icons */

  color: #5f6368;
  margin-bottom: 20px;
`;
