import React from 'react';
import styled from 'styled-components';
import { Button, Modal } from '../../components/common';
import QR from './QR';
import QrBar from './QrBar';
import { makeConnectQR } from '../../helpers';
import { colors } from '../../config/colors';

export default function ConnectCard(props) {
  const color = colors['light'];
  const { visible } = props;
  const { person } = props;

  const qrString = makeConnectQR(person?.owner_pubkey);

  return (
    <div onClick={(e) => e.stopPropagation()}>
      <Modal
        style={props.modalStyle}
        overlayClick={() => {
          props.dismiss();
        }}
        visible={visible}
      >
        <div style={{ textAlign: 'center', paddingTop: 59, width: 310 }}>
          <ImgWrap>
            <W color={color}>
              <Icon src={person?.img || '/static/person_placeholder.png'} />
            </W>
          </ImgWrap>
          <div
            style={{ textAlign: 'center', width: '100%', overflow: 'hidden', padding: '0 50px' }}
          >
            <N color={color}>Discuss this bounty with</N>
            <D color={color}>
              <B>{person?.owner_alias} </B>
            </D>

            <QR value={qrString} size={210} type={'connect'} />

            <QrBar value={person?.owner_pubkey} simple style={{ marginTop: 11 }} />

            <a href={qrString}>
              <Button
                text={'Connect with Sphinx'}
                color={'primary'}
                style={{ paddingLeft: 25, margin: '12px 0 40px' }}
                img={'sphinx_white.png'}
                imgSize={27}
                height={48}
                width={'100%'}
              />
            </a>
          </div>
        </div>
        <ModalBottomText color={color}>
          <img src="/static/scan_qr.svg" alt="scan" />
          <div className="bottomText">Scan or paste in Sphinx</div>
        </ModalBottomText>
      </Modal>
    </div>
  );
}

interface styledProps {
  color?: any;
}

const ImgWrap = styled.div`
  position: absolute;
  top: -45px;
  left: 0px;
  width: 100%;
  display: flex;
  justify-content: center;
  align-items: center;
  pointer-events: none;
  user-select: none;
`;
const B = styled.span`
  font-weight: bold;
`;
const W = styled.div<styledProps>`
  display: flex;
  align-items: center;
  justify-content: center;
  background: ${(p) => p?.color && p?.color.pureWhite};
  width: 90px;
  height: 90px;
  border-radius: 80px;
`;
const N = styled.div<styledProps>`
  font-family: Barlow;
  font-style: normal;
  font-weight: 500;
  font-size: 17px;
  line-height: 26px;
  text-align: center;
  color: ${(p) => p?.color && p?.color.grayish.G100};
`;

const D = styled.div<styledProps>`
  font-family: Barlow;
  font-style: normal;
  font-size: 20px;
  line-height: 26px;
  text-align: center;
  color: ${(p) => p?.color && p?.color.grayish.G10};
  margin-bottom: 20px;
`;

interface IconProps {
  src: string;
}

const Icon = styled.div<IconProps>`
  background-image: ${(p) => `url(${p.src})`};
  width: 80px;
  height: 80px;
  background-position: center; /* Center the image */
  background-repeat: no-repeat; /* Do not repeat the image */
  background-size: cover; /* Resize the background image to cover the entire container */
  border-radius: 80px;
  overflow: hidden;
`;

const ModalBottomText = styled.div<styledProps>`
  position: absolute;
  bottom: -36px;
  width: 310;
  background-color: transparent;
  display: flex;
  justify-content: center;
  .bottomText {
    margin-left: 12px;
    color: ${(p) => p?.color && p?.color.pureWhite};
  }
`;
