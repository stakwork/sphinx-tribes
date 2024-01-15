import React from 'react';
import { BountyMetrics } from 'store/main';
import copy from '../header/icons/copy.svg';
import hunter from '../header/icons/hunter.svg';
import coin from '../header/icons/coin.svg';

import {
  Wrapper,
  Card,
  VerticaGrayLine,
  DivWrapper,
  StatusWrapper,
  LeadingText,
  Title,
  TitleBlue,
  TitleGreen,
  Subheading,
  TitleWrapper,
  CardGreen,
  UpperCardWrapper,
  BelowCardWrapper,
  VerticaGrayLineSecondary,
  HorizontalGrayLine,
  CardHunter
} from './StatisticsStyles';
import './StatStyles.css';
export interface MockHunterMetrics {
  hunters_total_paid: number;
  hunters_first_bounty_paid: number;
}

interface StatisticsProps {
  metrics: BountyMetrics | undefined;
  freezeHeaderRef?: React.MutableRefObject<HTMLElement | null>;
  mockHunter?: MockHunterMetrics;
}

export const Statistics = ({ freezeHeaderRef, metrics, mockHunter }: StatisticsProps) => (
  <>
    <Wrapper ref={freezeHeaderRef}>
      <Card>
        <UpperCardWrapper>
          <TitleWrapper>
            <img className="BountiesSvg" src={copy} alt="" width="16.508px" height="20px" />
            <LeadingText>Bounties</LeadingText>
          </TitleWrapper>
          <StatusWrapper>
            <Subheading marginTop="5px" marginLeft="0px">
              Completed
            </Subheading>
            <TitleBlue>{metrics?.bounties_paid_average}%</TitleBlue>
          </StatusWrapper>
        </UpperCardWrapper>
        <HorizontalGrayLine />
        <BelowCardWrapper>
          <DivWrapper>
            <div>
              <Title>{metrics?.bounties_posted}</Title>
              <Subheading width="80px" data-testid="total_bounties_posted">
                Total Posted
              </Subheading>
            </div>
            <VerticaGrayLine />
          </DivWrapper>
          <DivWrapper>
            <div>
              <Title>78</Title>
              <Subheading>Assigned</Subheading>
            </div>
            <VerticaGrayLineSecondary />
          </DivWrapper>
          <DivWrapper>
            <div>
              <Title>{metrics?.bounties_paid}</Title>
              <Subheading data-testid="total_bounties_paid">Paid</Subheading>
            </div>
          </DivWrapper>
        </BelowCardWrapper>
      </Card>
      <CardGreen>
        <UpperCardWrapper>
          <TitleWrapper>
            <img className="SatoshieSvg" src={coin} alt="" width="23px" height="17px" />
            <LeadingText>Satoshis</LeadingText>
          </TitleWrapper>
          <StatusWrapper>
            <Subheading marginTop="5px" marginLeft="0px" data-testid="total_satoshis_paid">
              Paid
            </Subheading>
            <TitleGreen>{metrics?.sats_paid_percentage}%</TitleGreen>
          </StatusWrapper>
        </UpperCardWrapper>
        <HorizontalGrayLine />
        <BelowCardWrapper>
          <DivWrapper>
            <div>
              <Title>{metrics?.sats_posted}</Title>
              <Subheading data-testid="total_satoshis_posted">Total Posted</Subheading>
            </div>
            <VerticaGrayLine />
          </DivWrapper>
          <DivWrapper>
            <div>
              <Title>{metrics?.sats_paid}</Title>
              <Subheading>Paid</Subheading>
            </div>
            <VerticaGrayLineSecondary />
          </DivWrapper>
          <DivWrapper>
            <div>
              <Title>3 Days</Title>
              <Subheading width="120px">Avg Time to Paid</Subheading>
            </div>
          </DivWrapper>
        </BelowCardWrapper>
      </CardGreen>
      <CardHunter>
        <UpperCardWrapper>
          <TitleWrapper>
            <img className="HunterSvg" src={hunter} alt="" width="25px" height="25px" />
            <LeadingText>Hunters</LeadingText>
          </TitleWrapper>
        </UpperCardWrapper>
        <HorizontalGrayLine />
        <BelowCardWrapper>
          <DivWrapper>
            <div>
              <Title>{mockHunter?.hunters_total_paid}</Title>
              <Subheading width="80px">Total Paid</Subheading>
            </div>
            <VerticaGrayLine />
          </DivWrapper>
          <DivWrapper>
            <div>
              <Title>{mockHunter?.hunters_first_bounty_paid}</Title>
              <Subheading width="120px">First Bounty Paid</Subheading>
            </div>
          </DivWrapper>
        </BelowCardWrapper>
      </CardHunter>
    </Wrapper>
  </>
);
