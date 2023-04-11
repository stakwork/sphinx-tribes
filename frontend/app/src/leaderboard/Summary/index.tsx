import { EuiStat, EuiText, EuiTextColor } from '@elastic/eui';
import { colors } from 'config';
import { DollarConverter, satToUsd } from 'helpers';
import React from 'react';
import styled from 'styled-components';

export const Summary = ({ sats, bounties }: { sats: number; bounties: number }) => (
  <SummaryContainer>
    <EuiStat
      className="stats"
      title={`${DollarConverter(sats)}`}
      titleSize="s"
      titleColor={colors.light.black500}
      description={<EuiTextColor color={colors.light.black500}>SATS</EuiTextColor>}
    />
    <EuiStat
      title={DollarConverter(satToUsd(sats))}
      className="stats"
      titleSize="s"
      titleColor={colors.light.black500}
      description={<EuiTextColor color={colors.light.black500}>USD</EuiTextColor>}
    />
  </SummaryContainer>
);

const SummaryContainer = styled.div`
  display: flex;
  align-items: center;
  gap: 2rem;
  margin: auto;
  & .stats {
    background-color: ${colors.light.background};
    padding: 1rem 1rem 0 1rem;
    border: 1px solid ${colors.light.borderGreen1};
    border-radius: 0.5rem;
  }
`;
