import '@testing-library/jest-dom';
import { render, screen, waitFor } from '@testing-library/react';
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

const mockMetrics = {
  bounties_posted: 100,
  bounties_paid: 50,
  sats_posted: 5000,
  sats_paid: 2500,
  bounties_paid_average: 75,
  sats_paid_percentage: 50,
  average_paid: 10,
  average_completed: 1
};

describe('Statistics Component', () => {
  nock(user.url).get('/person/id/1').reply(200, {});

  test('display about view with extras', () => {
    render(<Statistics metrics={mockMetrics} />);
  });
});
