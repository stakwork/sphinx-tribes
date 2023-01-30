import React, { useRef, useState, useLayoutEffect } from 'react';
import styled from 'styled-components';
import { formatPrice, satToUsd } from '../../../helpers';
import { useIsMobile } from '../../../hooks';
import { Divider, Title, Paragraph } from '../../../sphinxUI';
import GalleryViewer from '../../utils/galleryViewer';
import NameTag from '../../utils/nameTag';
import FavoriteButton from '../../utils/favoriteButton';
import ReactMarkdown from 'react-markdown';

export function renderMarkdown(str) {
  return <ReactMarkdown>{str}</ReactMarkdown>;
}

export default function OfferSummary(props: any) {
  const { gallery, title, description, price, person, created, type } = props;

  const showPrice = !(type === 'offer_skill' || type === 'offer_other');

  const isMobile = useIsMobile();
  const [envHeight, setEnvHeight] = useState('100%');
  const imgRef: any = useRef(null);

  useLayoutEffect(() => {
    if (imgRef && imgRef.current) {
      if (imgRef.current?.offsetHeight > 100) {
        setEnvHeight(imgRef.current?.offsetHeight);
      }
    }
  }, [imgRef]);

  const heart = <FavoriteButton />;

  const nametag = (
    <NameTag
      {...person}
      iconSize={24}
      textSize={13}
      style={{ marginBottom: 14 }}
      created={created}
      widget={'offer'}
    />
  );

  if (isMobile) {
    return (
      <div style={{ padding: 20, overflow: 'auto' }}>
        <Pad>
          {nametag}

          <T>{title || 'No title'}</T>

          <Divider style={{ marginTop: 22 }} />
          <Y>
            {showPrice ? (
              <P>
                {formatPrice(price)} <B>SAT ({satToUsd(price)})</B>{' '}
              </P>
            ) : (
              <div />
            )}

            {heart}
          </Y>
          <Divider style={{ marginBottom: 22 }} />
          <D>{renderMarkdown(description)}</D>
          <GalleryViewer
            gallery={gallery}
            showAll={true}
            selectable={false}
            wrap={false}
            big={true}
          />
        </Pad>
      </div>
    );
  }

  return (
    <Wrap>
      <GalleryViewer
        innerRef={imgRef}
        style={{ width: 507, height: 'fit-content' }}
        gallery={gallery}
        showAll={false}
        selectable={false}
        wrap={false}
        big={true}
      />
      <div style={{ width: 316, padding: '40px 20px', overflowY: 'auto', height: envHeight }}>
        <Pad>
          {nametag}
          <Divider style={{ margin: '20px 0 20px' }} />
          <Title>{title}</Title>

          <Divider style={{ marginTop: 22 }} />
          <Y>
            {showPrice ? (
              <P>
                {formatPrice(price)} <B>SAT ({satToUsd(price)})</B>
              </P>
            ) : (
              <div />
            )}
            {heart}
          </Y>
          <Divider style={{ marginBottom: 22 }} />

          <Paragraph>{renderMarkdown(description)}</Paragraph>
        </Pad>
      </div>
    </Wrap>
  );
}
const Wrap = styled.div`
  display: flex;
  width: 100%;
  height: 100%;
  min-width: 800px;
  font-style: normal;
  font-weight: 500;
  font-size: 24px;
  line-height: 20px;
  color: #3c3f41;
  justify-content: space-between;
`;
const Pad = styled.div`
  padding: 0 20px;
`;
const Y = styled.div`
  display: flex;
  justify-content: space-between;
  width: 100%;
  height: 50px;
  align-items: center;
`;
const T = styled.div`
  font-weight: bold;
  font-size: 20px;
  margin: 10px 0;
`;
const B = styled.span`
  font-weight: 300;
`;
const P = styled.div`
  font-weight: 500;
`;
const D = styled.div`
  color: #5f6368;
  margin: 10px 0 30px;
`;
