import React, { useRef, useState } from 'react';
import styled from 'styled-components';
import { useObserver } from 'mobx-react-lite';
import { Button, Divider, Modal } from '../sphinxUI/index';

export default function Bot(props: any) {
  const {
    name,
    hideActions,
    small,
    id,
    img,
    tags,
    description,
    selected,
    select,
    created,
    owner_alias,
    owner_pubkey,
    unique_name,
    price_to_meet,
    extras,
    twitter_confirmed
  } = props;

  const defaultPic = '/static/bot_placeholder.png';
  const mediumPic = img;

  return useObserver(() => {
    function renderBotCard() {
      if (small) {
        return (
          <Wrap
            onClick={() => select(id, unique_name)}
            style={{
              background: selected ? '#F2F3F5' : '#fff'
            }}
          >
            <div>
              <Img src={mediumPic || defaultPic} style={hideActions && { width: 56, height: 56 }} />
            </div>
            <R style={{ width: hideActions ? 'calc(100% - 80px)' : 'calc(100% - 116px)' }}>
              <Title style={hideActions && { fontSize: 17 }}>{name}</Title>
              <Description>{description}</Description>
              {!hideActions && (
                <Row style={{ justifyContent: 'space-between', alignItems: 'center' }}>
                  <div />
                  <div style={{ height: 30 }} />
                </Row>
              )}
              <Divider style={{ marginTop: 20 }} />
            </R>
          </Wrap>
        );
      }
      // desktop mode
      return (
        <DWrap onClick={() => select(id, unique_name)}>
          <div>
            <Img
              style={{ height: 210, width: '100%', borderRadius: 0 }}
              src={mediumPic || defaultPic}
            />
            <div style={{ padding: 10 }}>
              <DTitle>{name}</DTitle>
              <DDescription>{description}</DDescription>
            </div>
          </div>
        </DWrap>
      );
    }

    return (
      <>
        {renderBotCard()}

        {/* <ConnectCard
                    dismiss={() => setShowQR(false)}
                    modalStyle={{ top: -64, height: 'calc(100% + 64px)' }}
                    person={props} visible={showQR} /> */}
      </>
    );
  });
}

const Wrap = styled.div`
  cursor: pointer;
  padding: 25px;
  padding-bottom: 0px;
  display: flex;
  width: 100%;
  overflow: hidden;
`;
const DWrap = styled.div`
  cursor: pointer;
  height: 320px;
  width: 210px;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  background: #fff;
  margin-bottom: 20px;
  margin-right: 20px;
  box-shadow: 0px 1px 2px rgba(0, 0, 0, 0.15);
  border-radius: 4px;
  overflow: hidden;
`;

const R = styled.div`
  margin-left: 20px;
`;

const Row = styled.div`
  display: flex;
  width: 100%;
`;

const Title = styled.h3`
  font-weight: 600;
  font-size: 17px;
  line-height: 19px;

  color: #3c3f41;

  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;

  /* Text 2 */

  color: #3c3f41;
`;

const DTitle = styled.h3`
  font-weight: 600;
  font-size: 17px;
  line-height: 19px;
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;

  color: #3c3f41;
`;
const Description = styled.div`
  font-size: 15px;
  color: #5f6368;
  white-space: nowrap;
  height: 26px;
  text-overflow: ellipsis;
  overflow: hidden;
  margin-bottom: 10px;
`;

const DDescription = styled.div`
  font-size: 12px;
  line-height: 18px;
  color: #5f6368;
  // white-space: nowrap;
  height: 36px;
  // text-overflow: ellipsis;
  overflow: hidden;
  // margin-bottom:10px;
`;
interface ImageProps {
  readonly src: string;
}
const Img = styled.div<ImageProps>`
  background-image: url('${(p) => p.src}');
  background-position: center;
  background-size: cover;
  height: 96px;
  width: 96px;
  position: relative;
  border-radius: 8px;
`;
