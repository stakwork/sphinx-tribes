import React from 'react';
import { BountyMetrics } from 'store/main';
import { convertToLocaleString } from 'helpers/helpers-extended';
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
// import './StatStyles.css';
export interface MockHunterMetrics {
  hunters_total_paid: number;
  hunters_first_bounty_paid: number;
}

interface StatisticsProps {
  metrics: BountyMetrics | undefined;
  freezeHeaderRef?: React.MutableRefObject<HTMLElement | null>;
  mockHunter?: MockHunterMetrics;
}

export const Statistics = ({ freezeHeaderRef, metrics }: StatisticsProps) => (
  <>
    <Wrapper ref={freezeHeaderRef}>
      <Card>
        <UpperCardWrapper>
          <TitleWrapper>
            <img
              style={{ marginTop: '4px', marginRight: '10px' }}
              src={copy}
              alt=""
              width="16.508px"
              height="20px"
            />
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
              <Title>{convertToLocaleString(metrics?.bounties_posted || 0)}</Title>
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
              <Title>{convertToLocaleString(metrics?.bounties_paid || 0)}</Title>
              <Subheading data-testid="total_bounties_paid">Paid</Subheading>
            </div>
          </DivWrapper>
        </BelowCardWrapper>
      </Card>
      <CardGreen>
        <UpperCardWrapper>
          <TitleWrapper>
            <img
              style={{ marginTop: '4px', marginRight: '10px' }}
              src={coin}
              alt=""
              width="23px"
              height="17px"
            />
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
              <Title>{convertToLocaleString(metrics?.sats_posted || 0)}</Title>
              <Subheading data-testid="total_satoshis_posted">Total Posted</Subheading>
            </div>
            <VerticaGrayLine />
          </DivWrapper>
          <DivWrapper>
            <div>
              <Title>{convertToLocaleString(metrics?.sats_paid || 0)}</Title>
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
            <img
              style={{ marginTop: '2px', marginRight: '10px' }}
              src={hunter}
              alt=""
              width="25px"
              height="25px"
            />
            <LeadingText>Hunters</LeadingText>
          </TitleWrapper>
        </UpperCardWrapper>
        <HorizontalGrayLine />
        <BelowCardWrapper>
          <DivWrapper>
            <div>
              <Title>{convertToLocaleString(metrics?.unique_hunters_paid || 0)}</Title>
              <Subheading width="80px">Total Paid</Subheading>
            </div>
            <VerticaGrayLine />
          </DivWrapper>
          <DivWrapper>
            <div>
              <Title>{convertToLocaleString(metrics?.new_hunters_paid || 0)}</Title>
              <Subheading width="120px">First Bounty Paid</Subheading>
            </div>
          </DivWrapper>
        </BelowCardWrapper>
      </CardHunter>
    </Wrapper>
  </>
);
