import { EuiText } from '@elastic/eui';
import { ImageButtonProps } from 'components/interfaces';
import React from 'react';
import styled from 'styled-components';

interface ButtonContainerProps {
  topMargin?: string;
  disabled?: boolean;
  paddingLeft?: string;
  paddingRight?: string;
}

const ButtonContainer = styled.div<ButtonContainerProps>`
  width: 220px;
  height: 48px;
  display: flex;
  align-items: center;
  justify-content: center;
  pointer-events: ${({ disabled }: ButtonContainerProps) => (disabled ? 'none' : 'all')};
  cursor: pointer;
  opacity: ${({ disabled }: ButtonContainerProps) => (disabled ? 0.8 : 1)};
  margin-top: ${(p: any) => p?.topMargin};
  background: #ffffff;
  border: 1px solid #dde1e5;
  border-radius: 30px;
  user-select: none;
  .ImageContainer {
    min-height: 48px;
    min-width: 48px;
    right: 37px;
    display: flex;
    align-items: center;
    justify-content: center;
  }
  .leadingImageContainer {
    padding-left: 20px;
    padding-right: 15px;
    display: flex;
    align-items: center;
    justify-content: center;
  }
  .buttonImage {
    filter: brightness(0) saturate(100%) invert(85%) sepia(10%) saturate(180%) hue-rotate(162deg)
      brightness(87%) contrast(83%);
  }
  :hover {
    border: 1px solid #b0b7bc;
  }
  :active {
    border: 1px solid #8e969c;
    .buttonImage {
      filter: brightness(0) saturate(100%) invert(22%) sepia(5%) saturate(563%) hue-rotate(161deg)
        brightness(91%) contrast(86%);
    }
  }
  .ButtonText {
    font-family: 'Barlow';
    font-style: normal;
    font-weight: 500;
    font-size: 14px;
    line-height: 17px;
    color: #5f6368;
  }
`;
const ImageButton = (props: ImageButtonProps) => (
  <ButtonContainer
    disabled={props.disabled}
    onClick={props?.buttonAction}
    style={{
      ...props.ButtonContainerStyle
    }}
  >
    {props.leadingImageSrc && (
      <div
        className="leadingImageContainer"
        style={{
          ...props.leadingImageContainerStyle
        }}
      >
        <img
          className="buttonImage"
          src={props.leadingImageSrc}
          alt={''}
          height={'14px'}
          width={'14px'}
        />
      </div>
    )}
    <EuiText
      className="ButtonText"
      style={{
        ...props.buttonTextStyle
      }}
    >
      {props.buttonText}
    </EuiText>
    {props.endImageSrc && (
      <div
        className="ImageContainer"
        style={{
          ...props.endingImageContainerStyle
        }}
      >
        <img
          className="buttonImage"
          src={props.endImageSrc}
          alt={'button_end_icon'}
          height={'12px'}
          width={'12px'}
        />
      </div>
    )}
  </ButtonContainer>
);

export default ImageButton;
