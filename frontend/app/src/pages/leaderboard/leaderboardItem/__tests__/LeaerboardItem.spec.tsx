import '@testing-library/jest-dom';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import nock from 'nock';
import React from 'react';
import { setupStore } from '../../../../__test__/__mockData__/setupStore';
import { user } from '../../../../__test__/__mockData__/user';
import { mockUsehistory } from '../../../../__test__/__mockFn__/useHistory';
import { LeaerboardItem } from '../index';

beforeAll(() => {
  nock.disableNetConnect();
  setupStore();
  mockUsehistory();
});

/**
 * @jest-environment jsdom
 */
describe('LeaerboardItem Component', () => {
  const owner_pubkey = '123456789';
  nock(user.url).get(`/person/id/1`).reply(200, {});
  nock(user.url).get(`/person/${owner_pubkey}`).reply(200, {});
  test('display leaderboard item with correct info', async () => {
    const mockOnClick = jest.fn();
    const convertedDollarAmount = '1 000';
    const props = {
      owner_pubkey: owner_pubkey,
      total_sats_earned: 1000,
      position: 1,
      total_bounties_completed: 1
    };

    const leaderboardItemRender = render(<LeaerboardItem key={0} {...props} />);

    //Expect converted dollar amount to be there
    expect(screen.queryByText(convertedDollarAmount)).toBeInTheDocument();
    expect(screen.queryByText('SAT')).toBeInTheDocument();
    expect(screen.queryByText(`#${props.position}`)).toBeInTheDocument();
    //expect(screen.queryByText(owner_pubkey)).toBeInTheDocument();
  });
});
