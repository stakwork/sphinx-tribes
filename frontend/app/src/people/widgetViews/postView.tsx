import React, { useState } from 'react';
import ReactMarkdown from 'react-markdown';
import styled from 'styled-components';
import { OfferViewProps } from 'people/interfaces';
import { Link, Paragraph, Title } from '../../components/common';
import GalleryViewer from '../utils/GalleryViewer';
import NameTag from '../utils/NameTag';

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

export default function PostView(props: OfferViewProps) {
  const { title, content, created, gallery, person } = props;
  const isLong = content && content.length > 100;
  const [expand, setExpand] = useState(false);

  const noGallery = !gallery || !gallery.length;

  return (
    <Wrap style={{ maxHeight: expand ? '' : 472 }}>
      <Pad>
        <Title>{title}</Title>
        <NameTag {...person} created={created} widget={'post'} />
        {/* // : <Date>{created && moment.unix(created).format('LLL')} </Date>} */}
        <Paragraph
          style={{
            maxHeight: noGallery && !expand ? 340 : !expand ? 80 : '',
            minHeight: noGallery ? 80 : '',
            overflow: 'hidden'
          }}
        >
          <ReactMarkdown>{content}</ReactMarkdown>
        </Paragraph>

        {isLong && (
          <Link
            style={{ textAlign: 'right', width: '100%' }}
            onClick={(e: any) => {
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
