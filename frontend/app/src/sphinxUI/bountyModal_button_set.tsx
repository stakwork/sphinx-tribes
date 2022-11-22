import { EuiText } from '@elastic/eui';
import React from 'react';
import styled from 'styled-components';
import { colors } from '../colors';

const ButtonSet = (props) => {
  const color = colors['light'];
  return (
    <ButtonSetContainer
      style={{
        ...props.ButtonSetContainerStyle
      }}
      color={color}>
      <ButtonContainer onClick={props?.githubShareAction} color={color}>
        <EuiText className="ButtonText">Github Ticket</EuiText>
        <div className="ImageContainer">
          <img
            className="buttonImage"
            src={'/static/github_ticket.svg'}
            alt={''}
            height={'14px'}
            width={'14px'}
          />
        </div>
      </ButtonContainer>
      <ButtonContainer topMargin={'16px'} onClick={props.copyURLAction} color={color}>
        <EuiText className="ButtonText">{props.copyStatus}</EuiText>
        <div className="ImageContainer">
          <img
            className="buttonImage"
            src={'/static/copy_link.svg'}
            alt={''}
            height={'10px'}
            width={'20px'}
          />
        </div>
      </ButtonContainer>
      <ButtonContainer topMargin={'16px'} onClick={props.twitterAction} color={color}>
        <EuiText className="ButtonText">Share to Twitter</EuiText>
        <div className="ImageContainer">
          <img
            className="buttonImage"
            src={'/static/share_with_twitter.svg'}
            alt={''}
            height={'15px'}
            width={'19px'}
          />
        </div>
      </ButtonContainer>
    </ButtonSetContainer>
  );
};

export default ButtonSet;

interface styledColor {
  color?: any;
}

interface ButtonContainerProps extends styledColor {
  topMargin?: string;
}

const ButtonSetContainer = styled.div`
  display: flex;
  flex-direction: column;
  padding-left: 36px;
  padding-top: 39px;
`;

const ButtonContainer = styled.div<ButtonContainerProps>`
  width: 220px;
  height: 48px;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-top: ${(p) => p?.topMargin};
  background: ${(p) => p?.color && p?.color.pureWhite};
  border: 1px solid ${(p) => p?.color && p?.color.grayish.G600};
  border-radius: 30px;
  user-select: none;
  .ImageContainer {
    position: absolute;
    min-height: 48px;
    min-width: 48px;
    right: 37px;
    display: flex;
    align-items: center;
    justify-content: center;
  }
  .buttonImage {
    filter: brightness(0) saturate(100%) invert(85%) sepia(10%) saturate(180%) hue-rotate(162deg)
      brightness(87%) contrast(83%);
  }
  :hover {
    border: 1px solid ${(p) => p?.color && p?.color.grayish.G300};
  }
  :active {
    border: 1px solid ${(p) => p?.color && p?.color.grayish.G100};
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
    color: ${(p) => p?.color && p?.color.grayish.G50};
  }
`;
