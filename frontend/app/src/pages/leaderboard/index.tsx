import { EuiLoadingSpinner } from '@elastic/eui';
import { colors, mobileBreakpoint } from 'config';
import { observer } from 'mobx-react-lite';
import React, { useEffect } from 'react';
import { useStores } from 'store';
import styled from 'styled-components';
import { LeaerboardItem } from './leaderboardItem';
import { Summary } from './summary';
import { Top3 } from './top3';

const Container = styled.div`
  height: calc(100% - 4rem);
  padding: 2rem;
  background-color: ${colors.light.background100};
  overflow: auto;
  align-items: center;
  justify-content: center;
  min-width: 400px;
  & > .inner {
    position: relative;
    margin: auto;
    max-width: 60%;
    min-width: 800px;
    display: flex;
    flex-direction: column;
    gap: 1rem;
  }

  & .summary {
    position: absolute;
    right: 0;
    top: 0;
  }

  @media (${mobileBreakpoint}) {
    height: calc(100% - 2rem);
    padding: 1rem;
    & > .inner {
      max-width: 100%;
      min-width: 300px;
    }
    & .summary {
      position: relative;
      right: 0;
      top: 0;
    }
  }
`;

const LoaderContainer = styled.div`
  display: flex;
  justify-content: center;
  align-items: center;
  width: 100%;
  height: 100%;
`;

export const LeaderboardPage = observer(() => {
  const { leaderboard } = useStores();
  useEffect(() => {
    leaderboard.fetchLeaders();
  }, [leaderboard]);

  if (leaderboard.isLoading) {
    return (
      <LoaderContainer>
        <EuiLoadingSpinner size="xl" />
      </LoaderContainer>
    );
  }
  return (
    <Container>
      <div className="inner">
        {leaderboard.total && (
          <Summary
            className="summary"
            bounties={leaderboard.total?.total_bounties_completed}
            sats={leaderboard.total?.total_sats_earned}
          />
        )}
        <Top3 />
        {leaderboard?.others.map((item: any, index: number) => (
          <LeaerboardItem position={index + 4} key={item.owner_pubkey} {...item} />
        ))}
      </div>
    </Container>
  );
});
