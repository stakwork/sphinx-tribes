import { EuiText } from '@elastic/eui';
import React from 'react';
import styled from 'styled-components';
import { ButtonContainer } from 'components/common';
import { colors } from '../config/colors';

const ButtonSetContainer = styled.div`
  display: flex;
  flex-direction: column;
  padding-left: 36px;
  padding-top: 39px;
  min-height: 300px;
`;

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
      {props.tribe !== 'None' ? (
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
      ) : (
        <ButtonContainer
          topMargin={'16px'}
          color={color}
          style={{ pointerEvents: 'none', opacity: 0.5 }}
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
            {props.tribe
              ? props.tribe.slice(0, 14) + (props.tribe.length > 14 ? '...' : '')
              : 'No Tribe'}
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
