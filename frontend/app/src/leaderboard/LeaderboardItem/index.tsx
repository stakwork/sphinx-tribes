import { EuiText } from '@elastic/eui';
import { PriceOuterContainer } from 'components/common';
import { colors } from 'config';
import { DollarConverter, satToUsd } from 'helpers';
import { UserInfo } from 'leaderboard/UserInfo';
import { LeaderItem } from 'leaderboard/store';

import React from 'react';
import styled from 'styled-components';

type Props = LeaderItem & {
  position: number;
};

const color = colors.light;
export const LeaerboardItem = ({
  owner_pubkey,
  total_bounties_completed,
  total_sats_earned,
  position
}: Props) => (
  <ItemContainer>
    <EuiText color={colors.light.text2} className="position">
      #{position}
    </EuiText>
    <UserInfo id={owner_pubkey} />
    <div className="userSummary">
      {/* <div className="bounties">
         Completed: {total_bounties_completed}
        </div> */}
      <div className="sats">
        <PriceOuterContainer
          price_Text_Color={color.primaryColor.P300}
          priceBackground={color.primaryColor.P100}
        >
          <div className="Price_inner_Container">
            <EuiText className="Price_Dynamic_Text">{DollarConverter(total_sats_earned)}</EuiText>
          </div>
          <div className="Price_SAT_Container">
            <EuiText className="Price_SAT_Text">SAT</EuiText>
          </div>
        </PriceOuterContainer>

        <EuiText color={color.grayish.G200} className="USD_Price">
          {DollarConverter(satToUsd(total_sats_earned))} <span className="currency">USD</span>
        </EuiText>
      </div>
    </div>
  </ItemContainer>
);

const ItemContainer = styled.div`
  position: relative;
  display: flex;
  flex-direction: row;
  align-items: center;
  gap: 1rem;
  padding: 1rem;
  background-color: ${colors.light.pureWhite};
  border-radius: 1rem;
  overflow: hidden;
  border: 1px solid transparent;
  transition-property: border box-shadow;
  transition-timing-function: ease;
  transition-duration: 0.2s;
  &:hover {
    border: 1px solid ${colors.light.borderGreen1};
    box-shadow: 0 0 5px 1px ${colors.light.borderGreen2};
  }

  & .userSummary {
    margin-left: auto;
    display: flex;
    align-items: center;
    gap: 1rem;
  }

  & .USD_Price {
    font-size: 1rem;
    text-align: right;
    .currency {
      font-size: 0.8em;
    }
  }
`;
