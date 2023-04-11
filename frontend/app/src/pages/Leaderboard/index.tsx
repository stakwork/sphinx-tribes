import { EuiLoadingSpinner } from '@elastic/eui';
import { colors } from 'config';
import { LeaerboardItem } from 'leaderboard/LeaderboardItem';
import { Summary } from 'leaderboard/Summary';
import { observer } from 'mobx-react-lite';
import React, { useEffect } from 'react';
import { useStores } from 'store';
import styled from 'styled-components';

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
            bounties={leaderboard.total?.total_bounties_completed}
            sats={leaderboard.total?.total_sats_earned}
          />
        )}
        {leaderboard?.sortedBySats.map((item, index) => (
          <LeaerboardItem position={index + 1} key={item.owner_pubkey} {...item} />
        ))}
      </div>
    </Container>
  );
});

const Container = styled.div`
  height: calc(100% - 4rem);
  padding: 2rem;
  background-color: ${colors.light.background100};
  overflow: auto;
  align-items: center;
  justify-content: center;
  min-width: 400px;
  & > .inner {
    margin: auto;
    max-width: 60%;
    min-width: 800px;
    display: flex;
    flex-direction: column;
    gap: 1rem;
  }

  @media (max-width: 900px) {
    height: calc(100% - 2rem);
    padding: 1rem;
    & > .inner {
      max-width: 100%;
      min-width: 300px;
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
