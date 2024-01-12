import '@testing-library/jest-dom';
import { render, screen, waitFor } from '@testing-library/react';
import nock from 'nock';
import React from 'react';
import { setupStore } from '../../../../__test__/__mockData__/setupStore';
import { user } from '../../../../__test__/__mockData__/user';
import { mockUsehistory } from '../../../../__test__/__mockFn__/useHistory';
import { Statistics } from '../';

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
  bounties_paid_average: 78,
  sats_paid_percentage: 50,
  average_paid: 10,
  average_completed: 1
};

describe('Statistics Component', () => {
  nock(user.url).get('/person/id/1').reply(200, {});

  it('renders without crashing', () => {
    const { container } = render(<Statistics metrics={mockMetrics} />);
    expect(container.firstChild).toBeInTheDocument();
  });

  it('renders bounties metrics correctly', () => {
    const { getByText } = render(<Statistics metrics={mockMetrics} />);
    expect(getByText('Bounties')).toBeInTheDocument();
    expect(getByText('100')).toBeInTheDocument();
    expect(getByText('Total Bounties Posted')).toBeInTheDocument();
    expect(getByText('50')).toBeInTheDocument();
    expect(getByText('Bounties Assigned')).toBeInTheDocument();
    expect(getByText('78')).toBeInTheDocument();
    expect(getByText('Bounties Paid')).toBeInTheDocument();
    expect(getByText('Completed')).toBeInTheDocument();
  });

  it('renders satoshis metrics correctly', () => {
    const { getByText } = render(<Statistics metrics={mockMetrics} />);
    expect(getByText('Satoshis')).toBeInTheDocument();
    expect(getByText('Total Sats Posted')).toBeInTheDocument();
    expect(getByText('Sats Paid')).toBeInTheDocument();
    expect(getByText('Avg Time to Paid')).toBeInTheDocument();
    expect(getByText('Paid')).toBeInTheDocument();
    expect(getByText('5000')).toBeInTheDocument();
    expect(getByText('2500')).toBeInTheDocument();
  });
});
