import '@testing-library/jest-dom';
import { render, screen } from '@testing-library/react';
import nock from 'nock';
import React from 'react';
import { setupStore } from '../../../__test__/__mockData__/setupStore';
import { user } from '../../../__test__/__mockData__/user';
import { mockUsehistory } from '../../../__test__/__mockFn__/useHistory';
import { Statistics } from './';

beforeAll(() => {
  nock.disableNetConnect();
  setupStore();
  mockUsehistory();
});

/**
 * @jest-environment jsdom
 */
describe('Statistics Component', () => {
  nock(user.url).get('/person/id/1').reply(200, {});

  test('display about view with extras', () => {
    const hardcodedWords = [
      'Bounties',
      '200',
      'Total Bounties Posted',
      '78',
      'Bounties Assigned',
      '136',
      '100%',
      '22536',
      'Total Sats Posted',
      '13625',
      'Sats Paid',
      '3 Days',
      'Avg Time to Paid',
      '48%',
      'Paid'
    ];

    const mockMetrics = {
      bounties_posted: 200,
      bounties_paid: 78,
      bounties_paid_average: 90,
      sats_posted: 1,
      sats_paid: 1,
      sats_paid_percentage: 1,
      average_paid: 1,
      average_completed: 1
    };

    render(<Statistics metrics={mockMetrics} />);

    const bountiesPaidElement = screen.getByText('Bounties Paid').nextSibling;

    expect(bountiesPaidElement).toHaveTextContent('78');

    for (let i = 0; i < hardcodedWords.length; i++) {
      expect(screen.getByText(hardcodedWords[i])).toBeInTheDocument();
    }
  });
});
