import React from 'react';
import AuthQR from 'people/auth/AuthQR';
import qr from './assets/qr.svg';
import exit from './assets/exit.svg';
import image2 from './assets/image 2.png';
import {
  FlexContainer,
  LeadingText,
  SecondaryText,
  TextWrapper,
  InfoText,
  ImageContainer,
  LogoSpan,
  ButtonDiv,
  ButtonText
} from './Styles/AuthStyles';

import './Styles/main.css';

export const Auth = () => (
  <FlexContainer>
    <TextWrapper>
      <LeadingText>Bounties</LeadingText>
      <SecondaryText>SuperAdmin</SecondaryText>
    </TextWrapper>
    <InfoText>Use Sphinx to sign in</InfoText>
    {/* <ImageContainer>
                <img src={image2} alt="" />
                <LogoSpan>
                   <img src={qr} alt="" />
                </LogoSpan>
            </ImageContainer> */}
    <AuthQR />
    <ButtonDiv onClick={console.log('clicked')}>
      Sign in with Sphinx
      <img className="exit" src={exit} alt="" />
    </ButtonDiv>
  </FlexContainer>
);
