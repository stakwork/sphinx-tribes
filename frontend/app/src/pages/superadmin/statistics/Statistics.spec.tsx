import '@testing-library/jest-dom';
import { queryByText, render, screen } from '@testing-library/react';
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

    render(<Statistics />);
     const bountiesPaid = screen.getAllByText('Bounties Paid'); 
     if(bountiesPaid){
      expect(bountiesPaid[0]).toBeInTheDocument();
     }
     for (let i = 0; i < hardcodedWords.length; i++) {
      expect(screen.queryByText(hardcodedWords[i])).toBeInTheDocument(); 
    } 
  });
});
