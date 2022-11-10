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
    bigClose
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
        ...style
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
