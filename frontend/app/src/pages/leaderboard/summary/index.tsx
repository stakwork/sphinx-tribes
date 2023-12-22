import { EuiStat, EuiTextColor } from '@elastic/eui';
import { colors, mobileBreakpoint } from 'config';
import { DollarConverter } from 'helpers';
import React from 'react';
import styled from 'styled-components';

const SummaryContainer = styled.div`
  display: flex;
  flex-direction: column;
  align-items: stretch;
  gap: 1rem;
  margin: auto;
  & .stats {
    background-color: ${colors.light.background};
    padding: 1rem 1rem 0 1rem;
    border: 1px solid ${colors.light.borderGreen1};
    border-radius: 0.5rem;
  }
  @media (${mobileBreakpoint}) {
    flex-direction: row;
  }
`;

export const Summary = ({
  sats,
  bounties,
  className = ''
}: {
  sats: number;
  bounties: number;
  className?: string;
}) => (
  <SummaryContainer className={className}>
    <EuiStat
      className="stats"
      title={`${DollarConverter(sats)}`}
      titleSize="s"
      titleColor={colors.light.black500}
      description={<EuiTextColor color={colors.light.black500}>Total sats earned</EuiTextColor>}
    />
    <EuiStat
      title={bounties}
      className="stats"
      titleSize="s"
      titleColor={colors.light.black500}
      description={<EuiTextColor color={colors.light.black500}>Total tasks completed</EuiTextColor>}
    />
  </SummaryContainer>
);
