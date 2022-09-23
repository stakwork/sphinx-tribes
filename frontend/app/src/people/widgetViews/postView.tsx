import React, { useState } from 'react';
import styled from 'styled-components';
import GalleryViewer from '../utils/galleryViewer';
import ReactMarkdown from 'react-markdown';
import NameTag from '../utils/nameTag';
import { Title, Paragraph, Date, Link } from '../../sphinxUI';

export default function PostView(props: any) {
  const { title, content, created, gallery, showName, person } = props;
  const isLong = content && content.length > 100;
  const [expand, setExpand] = useState(false);

  const noGallery = !gallery || !gallery.length;

  return (
    <Wrap style={{ maxHeight: expand ? '' : 472 }}>
      <Pad>
        <Title>{title} </Title>
        <NameTag {...person} created={created} widget={'post'} />
        {/* // : <Date>{created && moment.unix(created).format('LLL')} </Date>} */}
        <Paragraph
          style={{
            maxHeight: noGallery && !expand ? 340 : !expand ? 80 : '',
            minHeight: noGallery ? 80 : '',
            overflow: 'hidden'
          }}
        >
          <ReactMarkdown>{content}</ReactMarkdown>{' '}
        </Paragraph>

        {isLong && (
          <Link
            style={{ textAlign: 'right', width: '100%' }}
            onClick={(e) => {
              e.stopPropagation();
              setExpand(!expand);
            }}
          >
            {!expand ? 'Read more' : 'Show less'}
          </Link>
        )}
      </Pad>
      <GalleryViewer
        cover
        showAll={false}
        big={true}
        wrap={false}
        selectable={true}
        gallery={gallery}
        style={{ maxHeight: 291, overflow: 'hidden' }}
      />
    </Wrap>
  );
}

const Wrap = styled.div`
  display: flex;
  flex-direction: column;
  width: 100%;
  min-height: 472px;
  min-width: 100%;
  font-style: normal;
  font-weight: 500;
  font-size: 24px;
  line-height: 20px;
  color: #3c3f41;
  justify-content: space-between;
`;

const Pad = styled.div`
  display: flex;
  flex-direction: column;
  padding: 20px;
`;
