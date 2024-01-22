import React, { useState } from 'react';
import styled from 'styled-components';
import { getHost } from '../../config';
import { Button, Divider, LazyImgBg } from '../../components/common';
import ConnectCard from '../../people/utils/ConnectCard';
import { PersonProps } from '../../people/interfaces';

const Wrap = styled.div`
  cursor: pointer;
  padding: 10px;
  min-height: 88px;
  display: flex;
  align-items: center;
  justify-content: center;
  width: 100%;
  overflow: hidden;
`;

interface DWarpProps {
  squeeze: boolean;
}
const DWrap = styled.div<DWarpProps>`
  cursor: pointer;
  height: 350px;
  width: ${(p: any) => (p.squeeze ? '200px' : '210px')};
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  background: #fff;
  margin-bottom: 20px;
  margin-right: 20px;
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.15);
  border-radius: 4px;
  overflow: hidden;
`;

const R = styled.div`
  margin-left: 20px;
  display: flex;
  flex-direction: column;
  justify-content: center;
`;

const Row = styled.div`
  display: flex;
  width: 100%;
`;

const Title = styled.h3`
  font-weight: 500;
  font-size: 17px;
  /* Text 2 */

  color: #3c3f41;
`;

const DTitle = styled.h3`
  font-weight: 500;
  font-size: 17px;
  line-height: 19px;

  color: #3c3f41;
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
`;

interface DescriptionProps {
  lineRows: number;
  isMobile: boolean;
}
const Description = styled.div<DescriptionProps>`
  font-size: 12px;
  color: #5f6368;
  text-overflow: ellipsis;
  overflow: hidden;
  margin-bottom: 10px;
  font-weight: 400;
  max-width: ${(p: any) => (p.isMobile ? '200px' : 'auto')};

  display: -webkit-box;
  -webkit-line-clamp: ${(p: any) => (p.lineRows ? p.lineRows : 1)};
  -webkit-box-orient: vertical;
`;

const DDescription = styled.div`
  font-size: 12px;
  line-height: 18px;
  color: #5f6368;
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
`;

const Img = styled(LazyImgBg)`
  background-position: center;
  background-size: cover;
  height: 96px;
  width: 96px;
  border-radius: 50%;
  position: relative;
`;
const host = getHost();

function makeQR(pubkey: string) {
  return `sphinx.chat://?action=person&host=${host}&pubkey=${pubkey}`;
}

export default function Person(props: PersonProps) {
  const {
    hideActions,
    small,
    id,
    img,
    selected,
    select,
    owner_alias,
    owner_pubkey,
    unique_name,
    squeeze,
    description
  } = props;
  const defaultPic = '/static/person_placeholder.png';
  const addedStyles = hideActions ? { width: 56, height: 56 } : {};
  const qrString = makeQR(owner_pubkey);
  const [showQR, setShowQR] = useState(false);
  function renderPersonCard() {
    if (small) {
      return (
        <a
          href={`/p/${owner_pubkey}`}
          style={{ textDecoration: 'none', color: 'inherit', cursor: 'pointer' }}
        >
          <div
            style={{ background: selected ? '#F2F3F5' : '#fff', cursor: 'pointer' }}
            onClick={() => select(id, unique_name, owner_pubkey)}
          >
            <Wrap style={{ padding: hideActions ? 10 : 25 }}>
              <div>
                <Img style={addedStyles} src={img || defaultPic} />
              </div>
              <R style={{ width: hideActions ? 'calc(100% - 80px)' : 'calc(100% - 116px)' }}>
                <Title style={{ fontSize: hideActions ? 17 : 20, margin: 0 }}>{owner_alias}</Title>
                {description && description !== 'description' && (
                  <Description
                    lineRows={hideActions ? 1 : 2}
                    style={{
                      margin: 0,
                      marginTop: hideActions ? 5 : 10,
                      fontSize: hideActions ? 12 : 15
                    }}
                    isMobile={small}
                  >
                    {description}
                  </Description>
                )}
                {!hideActions && (
                  <Row style={{ justifyContent: 'space-between', alignItems: 'center' }}>
                    {!hideActions && owner_pubkey ? (
                      <>
                        <a href={qrString}>
                          <Button
                            text="Connect"
                            color="white"
                            leadingIcon={'open_in_new'}
                            iconSize={16}
                            style={{ marginTop: 12 }}
                            onClick={(e: any) => e.stopPropagation()}
                          />
                        </a>
                      </>
                    ) : (
                      <div style={{ height: 30 }} />
                    )}
                  </Row>
                )}
              </R>
            </Wrap>
            <Divider />
          </div>
        </a>
      );
    } else {
      // desktop mode
      return (
        <a
          href={`/p/${owner_pubkey}`}
          style={{ textDecoration: 'none', color: 'inherit', cursor: 'pointer' }}
        >
          <DWrap squeeze={squeeze} onClick={() => select(id, unique_name, owner_pubkey)}>
            <div>
              <div style={{ height: 210 }}>
                <Img style={{ height: '100%', width: '100%', borderRadius: 0 }} src={img} />
              </div>
              <div style={{ padding: 16 }}>
                <DTitle>{owner_alias}</DTitle>
                {description && description !== 'description' && (
                  <DDescription>{description}</DDescription>
                )}
              </div>
              <div>
                <Divider />
                <Row style={{ justifyContent: 'space-between', alignItems: 'center', height: 50 }}>
                  {owner_pubkey ? (
                    <>
                      <Button
                        text="Connect"
                        color="clear"
                        iconStyle={{ color: '#B0B7BC' }}
                        endingIcon={'open_in_new'}
                        style={{ fontSize: 13, fontWeight: 500 }}
                        iconSize={16}
                        onClick={(e: any) => {
                          setShowQR(true);
                          e.stopPropagation();
                        }}
                      />
                    </>
                  ) : (
                    <div />
                  )}
                </Row>
              </div>
            </div>
          </DWrap>
        </a>
      );
    }
  }

  return (
    <>
      {renderPersonCard()}

      <ConnectCard
        dismiss={() => setShowQR(false)}
        modalStyle={{ top: -64, height: 'calc(100% + 64px)' }}
        person={props}
        visible={showQR}
      />
    </>
  );
}
