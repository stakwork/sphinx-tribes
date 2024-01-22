import React from 'react';
import styled from 'styled-components';
import { ModalProps } from 'components/interfaces';
import FadeLeft from '../animated/FadeLeft';
import { colors } from '../../config/colors';
import { IconButton } from '.';
import { Portal } from './index';

interface styledProps {
  color?: any;
}

const R = styled.div<styledProps>`
  position: absolute;
  right: -75px;
  top: 0px;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
`;

const L = styled.div<styledProps>`
  position: absolute;
  left: -75px;
  top: 0px;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
`;

const RNew = styled.div<styledProps>`
  position: absolute;
  right: -63.9px;
  top: 0px;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
`;

const LNew = styled.div<styledProps>`
  position: absolute;
  left: -62.9px;
  top: 0px;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
`;

const CircL = styled.div<styledProps>`
  display: flex;
  align-items: center;
  justify-content: center;
  width: 62px;
  height: 88px;
  background: ${(p: any) => p.color && p.color.black150};
  border-radius: 10px 0px 0px 10px;
  cursor: pointer;
`;

const CircR = styled.div<styledProps>`
  display: flex;
  align-items: center;
  justify-content: center;
  width: 62px;
  height: 88px;
  background: ${(p: any) => p.color && p.color.black150};
  border-radius: 0px 10px 10px 0px;
  cursor: pointer;
`;

const Circ = styled.div<styledProps>`
  display: flex;
  align-items: center;
  justify-content: center;
  width: 65px;
  height: 65px;
  background: ${(p: any) => p.color && p.color.grayish.G60};
  border-radius: 50px;
  cursor: pointer;
`;

const X = styled.div<styledProps>`
  position: absolute;
  top: 5px;
  right: 8px;
  cursor: pointer;
`;

const BigX = styled.div<styledProps>`
  position: absolute;
  top: 20px;
  right: -58px;
  cursor: pointer;
  z-index: 10;
`;
const Env = styled.div<styledProps>`
  width: 312px;
  min-height: 254px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 16px;
  background: ${(p: any) => p.color && p.color.pureWhite};
  position: relative;
`;

export default function Modal(props: ModalProps) {
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
    bigCloseImage,
    bigCloseImageStyle
  } = props;

  const color = colors['light'];
  const fillStyle = fill
    ? {
        minHeight: '100%',
        height: '100%',
        width: '100%',
        borderRadius: 0
      }
    : {};

  return (
    <Portal>
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
          overflowY: 'auto',
          overflowX: 'hidden',
          ...style
        }}
      >
        <Env
  
          role="alertdialog"
          aria-modal={true}
          style={{ ...fillStyle, ...envStyle }}
          color={color}
        >
          {close && (
            <X color={color}>
              <IconButton onClick={close} size={20} icon="close" color="transparent" />
            </X>
          )}

          {bigClose && (
            <BigX color={color}>
              <IconButton onClick={bigClose} size={36} icon="close" color="transparent" />
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
                zIndex: 10,
                ...bigCloseImageStyle
              }}
              onClick={bigCloseImage}
            >
              <img src="/static/Close.svg" alt="close_svg" height={'100%'} width={'100%'} />
            </div>
          )}

          {prevArrow && (
            <L color={color}>
              <Circ color={color}>
                <IconButton
                  color="noColor"
                  iconStyle={{ color: color.pureWhite }}
                  icon={'chevron_left'}
                  onClick={(e: any) => {
                    e.stopPropagation();
                    prevArrow();
                  }}
                />
              </Circ>
            </L>
          )}
          {nextArrow && (
            <R color={color}>
              <Circ color={color}>
                <IconButton
                  color="noColor"
                  icon={'chevron_right'}
                  iconStyle={{ color: color.pureWhite }}
                  onClick={(e: any) => {
                    e.stopPropagation();
                    nextArrow();
                  }}
                />
              </Circ>
            </R>
          )}

          {prevArrowNew && (
            <LNew color={color}>
              <CircL color={color}>
                <IconButton
                  iconStyle={{ color: color.pureWhite }}
                  icon={'chevron_left'}
                  onClick={(e: any) => {
                    e.stopPropagation();
                    prevArrowNew();
                  }}
                />
              </CircL>
            </LNew>
          )}
          {nextArrowNew && (
            <RNew color={color}>
              <CircR color={color}>
                <IconButton
                  icon={'chevron_right'}
                  iconStyle={{ color: color.pureWhite }}
                  onClick={(e: any) => {
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
    </Portal>
  );
}
