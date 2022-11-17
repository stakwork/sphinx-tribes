import React from 'react';
import styled from 'styled-components';
import FadeLeft from '../animated/fadeLeft';
import { IconButton } from '.';

export default function Modal(props: any) {
  const {
    visible,
    fill,
    overlayClick,
    dismountCallback,
    children,
    close,
    style,
    hideOverlay,
    envStyle,
    nextArrow,
    prevArrow,
    nextArrowNew,
    prevArrowNew,
    bigClose,
    bigCloseImage
  } = props;

  const fillStyle = fill
    ? {
        height: '100%',
        width: '100%',
        borderRadius: 0
      }
    : {};

  return (
    <FadeLeft
      withOverlay={!hideOverlay}
      drift={100}
      direction="up"
      close={close || bigClose}
      overlayClick={overlayClick}
      dismountCallback={dismountCallback}
      isMounted={visible ? true : false}
      style={{
        ...style,
        position: 'absolute',
        top: 0,
        left: 0,
        zIndex: 1000000,
        width: '100%',
        height: '100%',
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
        // overflow: 'auto',
      }}
    >
      <Env style={{ ...fillStyle, ...envStyle }}>
        {close && (
          <X>
            <IconButton onClick={close} size={20} icon="close" />
          </X>
        )}

        {bigClose && (
          <BigX>
            <IconButton onClick={bigClose} size={36} icon="close" />
          </BigX>
        )}

        {bigCloseImage && (
          <div
            style={{
              height: '40px',
              width: '40px',
              position: 'absolute',
              top: '8px',
              right: '-48px',
              cursor: 'pointer',
              zIndex: 10
            }}
            onClick={bigCloseImage}>
            <img src="static/Close.svg" alt="close_svg" height={'100%'} width={'100%'} />
          </div>
        )}

        {prevArrow && (
          <L>
            <Circ>
              <IconButton
                iconStyle={{ color: '#fff' }}
                icon={'chevron_left'}
                onClick={(e) => {
                  e.stopPropagation();
                  prevArrow();
                }}
              />
            </Circ>
          </L>
        )}
        {nextArrow && (
          <R>
            <Circ>
              <IconButton
                icon={'chevron_right'}
                iconStyle={{ color: '#fff' }}
                onClick={(e) => {
                  e.stopPropagation();
                  nextArrow();
                }}
              />
            </Circ>
          </R>
        )}

        {prevArrowNew && (
          <LNew>
            <CircL>
              <IconButton
                iconStyle={{ color: '#fff' }}
                icon={'chevron_left'}
                onClick={(e) => {
                  e.stopPropagation();
                  prevArrowNew();
                }}
              />
            </CircL>
          </LNew>
        )}
        {nextArrowNew && (
          <RNew>
            <CircR>
              <IconButton
                icon={'chevron_right'}
                iconStyle={{ color: '#fff' }}
                onClick={(e) => {
                  e.stopPropagation();
                  nextArrowNew();
                }}
              />
            </CircR>
          </RNew>
        )}
        {children}
      </Env>
    </FadeLeft>
  );
}

const R = styled.div`
  position: absolute;
  right: -85px;
  top: 0px;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
`;

const L = styled.div`
  position: absolute;
  left: -85px;
  top: 0px;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
`;

const RNew = styled.div`
  position: absolute;
  right: -63.9px;
  top: 0px;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
`;

const LNew = styled.div`
  position: absolute;
  left: -62.9px;
  top: 0px;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
`;

const CircL = styled.div`
  display: flex;
  align-items: center;
  justify-content: center;
  width: 62px;
  height: 88px;
  background: rgba(0, 0, 0, 0.75);
  border-radius: 10px 0px 0px 10px;
  cursor: pointer;
`;

const CircR = styled.div`
  display: flex;
  align-items: center;
  justify-content: center;
  width: 62px;
  height: 88px;
  background: rgba(0, 0, 0, 0.75);
  border-radius: 0px 10px 10px 0px;
  cursor: pointer;
`;

const Circ = styled.div`
  display: flex;
  align-items: center;
  justify-content: center;
  width: 65px;
  height: 65px;
  background: #ffffff44;
  border-radius: 50px;
  cursor: pointer;
`;

const X = styled.div`
  position: absolute;
  top: 5px;
  right: 8px;
  cursor: pointer;
`;

const BigX = styled.div`
  position: absolute;
  top: 20px;
  right: -68px;
  cursor: pointer;
  z-index: 10;
`;
const Env = styled.div`
  width: 312px;
  min-height: 254px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 16px;
  background: #ffffff;
  position: relative;
`;
