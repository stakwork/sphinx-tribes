import React from 'react';
import copy from '../Header/icons/copy.svg';
import copygray from '../Header/icons/copygray.svg';
import bountiesposted from '../Header/icons/bountiesposted.svg';
import bountiespaid from '../Header/icons/bountiespaid.svg';
import clockloader from '../Header/icons/clockloder.svg';
import coin from '../Header/icons/coin.svg';
import satoshiesposted from '../Header/icons/satoshiesposted.svg';
import satoshiespaid from '../Header/icons/satoshiespaid.svg';
import calendar from '../Header/icons/calendar.svg';
import {
  Wrapper,
  Card,
  VerticaGrayLine,
  DivWrapper,
  LeadingText,
  Title,
  TitleBlue,
  TitleGreen,
  Subheading,
  TitleWrapper,
  CardGreen,
  VerticaGrayLineSecondary,
  VerticaGrayLineAleternative
} from './StatisticsStyles';
import './StatStyles.css';

export const Statistics = () => {
  console.log('super admin');
  return (
    <>
      <Wrapper>
        <Card>
          <TitleWrapper>
            <img className="BountiesSvg" src={copy} alt="" width="16.508px" height="20px" />
            <LeadingText>Bounties</LeadingText>
          </TitleWrapper>

          <DivWrapper>
            <VerticaGrayLine />
            <img className="logoAlign" src={bountiesposted} alt="" width="27.09px" height="20px" />
            <div>
              <Title>200</Title>
              <Subheading>Total Bounties Posted</Subheading>
            </div>
          </DivWrapper>
          <DivWrapper>
            <VerticaGrayLineAleternative />
            <img className="logoAlign" src={copygray} alt="" width="27.09px" height="20px" />
            <div>
              <Title>78</Title>
              <Subheading className="BounA">Bounties Assigned</Subheading>
            </div>
          </DivWrapper>
          <DivWrapper>
            <VerticaGrayLine />
            <img className="logoAlign" src={bountiespaid} alt="" width="20" height="20" />
            <div>
              <Title>136</Title>
              <Subheading>Bounties Paid</Subheading>
            </div>
          </DivWrapper>
          <DivWrapper>
            <VerticaGrayLineSecondary />
            <img className="ClocklogoAlign" src={clockloader} alt="" width="24px" height="24px" />
            <div>
              <TitleBlue>100%</TitleBlue>
              <Subheading>Bounties Paid</Subheading>
            </div>
          </DivWrapper>
        </Card>
        <CardGreen>
          <TitleWrapper>
            <img className="SatoshieSvg" src={coin} alt="" width="23px" height="17px" />
            <LeadingText>Satoshis</LeadingText>
          </TitleWrapper>
          <DivWrapper>
            <VerticaGrayLine />
            <img className="logoAlign" src={satoshiesposted} alt="" width="23px" height="17" />
            <div>
              <Title>22536</Title>
              <Subheading>Total Sats Posted</Subheading>
            </div>
          </DivWrapper>
          <DivWrapper>
            <VerticaGrayLineAleternative />
            <img className="logoAlign" src={satoshiespaid} alt="" width="23px" height="17px" />
            <div>
              <Title>13625</Title>
              <Subheading>Sats Paid</Subheading>
            </div>
          </DivWrapper>
          <DivWrapper>
            <VerticaGrayLine />
            <img className="logoAlignSecondary" src={calendar} alt="" width="23px" height="17px" />
            <div>
              <Title>3 Days</Title>
              <Subheading>Avg Time to Paid</Subheading>
            </div>
          </DivWrapper>
          <DivWrapper>
            <VerticaGrayLineSecondary />
            <img className="ClocklogoAlign" src={clockloader} alt="" width="24px" height="24px" />
            <div>
              <TitleGreen>48%</TitleGreen>
              <Subheading>Paid</Subheading>
            </div>
          </DivWrapper>
        </CardGreen>
      </Wrapper>
    </>
  );
};
