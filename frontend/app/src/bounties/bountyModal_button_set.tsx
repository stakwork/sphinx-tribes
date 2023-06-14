import { EuiText } from '@elastic/eui';
import React from 'react';
import styled from 'styled-components';
import { colors } from '../config/colors';

const ButtonSet = ({ showGithubBtn, ...props }: any) => {
  const color = colors['light'];
  return (
    <ButtonSetContainer
      style={{
        ...props.ButtonSetContainerStyle
      }}
    >
      {showGithubBtn && (
        <ButtonContainer onClick={props?.githubShareAction} color={color}>
          <div className="LeadingImageContainer">
            <img
              className="buttonImage"
              src={'/static/github_icon.svg'}
              alt={'github_ticket'}
              height={'20px'}
              width={'20px'}
            />
          </div>
          <EuiText className="ButtonText">Github Ticket</EuiText>
          <div className="ImageContainer">
            <img
              className="buttonImage"
              src={'/static/github_ticket.svg'}
              alt={'github_ticket'}
              height={'14px'}
              width={'14px'}
            />
          </div>
        </ButtonContainer>
      )}
      {props?.replitLink && (
        <ButtonContainer
          topMargin={'16px'}
          onClick={() => {
            window.open(props.replitLink[0]);
          }}
          color={color}
        >
          <div
            className="LeadingImageContainer"
            style={{
              marginLeft: '20px'
            }}
          >
            <img
              className="buttonImage"
              src={'/static/replit_icon.svg'}
              alt={'github_ticket'}
              height={'18px'}
              width={'14px'}
            />
          </div>
          <EuiText className="ButtonText">Replit</EuiText>
          <div className="ImageContainer">
            <img
              className="buttonImage"
              src={'/static/github_ticket.svg'}
              alt={'github_ticket'}
              height={'14px'}
              width={'14px'}
            />
          </div>
        </ButtonContainer>
      )}
      {props.tribe && (
        <ButtonContainer
          topMargin={'16px'}
          onClick={() => {
            props?.tribeFunction();
          }}
          color={color}
        >
          <div
            className="LeadingImageContainer"
            style={{
              marginLeft: '6px',
              marginRight: '12px'
            }}
          >
            <img
              src={'/static/tribe_demo.svg'}
              alt={'github_ticket'}
              height={'32px'}
              width={'32px'}
            />
          </div>
          <EuiText className="ButtonText">
            {props.tribe.slice(0, 14)} {props.tribe.length > 14 && '...'}
          </EuiText>
          <div className="ImageContainer">
            <img
              className="buttonImage"
              src={'/static/github_ticket.svg'}
              alt={'github_ticket'}
              height={'14px'}
              width={'14px'}
            />
          </div>
        </ButtonContainer>
      )}

      <ButtonContainer topMargin={'16px'} onClick={props.copyURLAction} color={color}>
        <div className="LeadingImageContainer">
          <img
            className="buttonImage"
            src={'/static/copy_icon_link.svg'}
            alt={'copy_link'}
            height={'20px'}
            width={'20px'}
          />
        </div>
        <EuiText className="ButtonText">{props.copyStatus}</EuiText>
      </ButtonContainer>
      <ButtonContainer topMargin={'16px'} onClick={props.twitterAction} color={color}>
        <div className="LeadingImageContainer">
          <img
            className="buttonImage"
            src={'/static/share_with_twitter.svg'}
            alt={'twitter'}
            height={'15px'}
            width={'19px'}
          />
        </div>
        <EuiText className="ButtonText">Share to Twitter</EuiText>
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
  min-height: 300px;
`;

const ButtonContainer = styled.div<ButtonContainerProps>`
  width: 220px;
  height: 48px;
  display: flex;
  align-items: center;
  justify-content: flex-start;
  margin-top: ${(p: any) => p?.topMargin};
  background: ${(p: any) => p?.color && p?.color.pureWhite};
  border: 1px solid ${(p: any) => p?.color && p?.color.grayish.G600};
  border-radius: 30px;
  user-select: none;
  .LeadingImageContainer {
    margin-left: 14px;
    margin-right: 16px;
  }
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
    border: 1px solid ${(p: any) => p?.color && p?.color.grayish.G300};
  }
  :active {
    border: 1px solid ${(p: any) => p?.color && p?.color.grayish.G100};
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
    color: ${(p: any) => p?.color && p?.color.grayish.G50};
  }
`;
