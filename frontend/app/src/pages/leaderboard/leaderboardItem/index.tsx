import { EuiText } from '@elastic/eui';
import React from 'react';
import styled from 'styled-components';
import { PriceOuterContainer } from '../../../components/common';
import { colors } from '../../../config';
import { DollarConverter } from '../../../helpers';
import { UserInfo } from '../userInfo';
import { LeaderItem } from '../store';

const ItemContainer = styled.div`
  --position-gutter: 3rem;
  position: relative;
  display: flex;
  flex-direction: row;
  align-items: center;
  gap: 1rem;
  padding: 0.5rem;
  margin-left: var(--position-gutter);
  background-color: ${colors.light.pureWhite};
  border-radius: 0.5rem;
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

  & .position {
    position: absolute;
    left: calc(-1 * var(--position-gutter));
    font-weight: 500;
  }
`;

type Props = LeaderItem & {
  position: number;
  owner_pubkey: string;
  total_sats_earned: number;
};

const color = colors.light;
export const LeaerboardItem = ({ owner_pubkey, total_sats_earned, position }: Props) => (
  <ItemContainer>
    <EuiText color={colors.light.text2} className="position">
      #{position}
    </EuiText>
    <UserInfo id={owner_pubkey} />
    <div className="userSummary">
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
      </div>
    </div>
  </ItemContainer>
);
